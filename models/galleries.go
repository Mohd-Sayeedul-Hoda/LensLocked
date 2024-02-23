package models

import (
	"github.com/jinzhu/gorm"
)

type Gallery struct{
  gorm.Model
  UserID uint `gorm:"not_null;index"`
  Title string `gorm:"not_null"`
}

type galleryGorm struct {
  db *gorm.DB
}

type galleryService struct{
  GalleryDB
}

type galleryValidator struct{
  GalleryDB
}

type GalleryDB interface {
  Create(gallery *Gallery)error
  ByID(id uint)(*Gallery, error)
  ByUserID(userID uint)([]Gallery, error)
  Update(gallery *Gallery) error
  Delete(id uint)error
}

type GalleryService interface {
  GalleryDB
}

var _ GalleryDB = &galleryGorm{}

type galleryValFn func(*Gallery)error

const(
  ErrUserIDRequired modelError = "models: user ID is required"
  ErrTitleRequire modelError = "models: title is required"
)

func NewGalleryService(db *gorm.DB) GalleryService{
  return &galleryService{
    GalleryDB: &galleryValidator{
      GalleryDB: &galleryGorm{
	db: db,
      },
    },
  }
}

func (gg *galleryGorm) Create(gallery *Gallery) error{
  return gg.db.Create(gallery).Error
}

func (gg *galleryGorm) ByID(id uint)(*Gallery, error){
  var gallery Gallery
  db := gg.db.Where("id = ?", id)
  err := first(db, &gallery)
  if err != nil{
    return nil, err
  }
  return &gallery, nil
}

func (gg *galleryGorm) Update(gallery *Gallery)error{
  return gg.db.Save(gallery).Error
}

func (gg *galleryGorm) Delete(id uint)error{
  gallery := Gallery{Model: gorm.Model{ID: id}}
  return gg.db.Delete(&gallery).Error
}

func (gg *galleryGorm) ByUserID(userID uint) ([]Gallery, error){
  var galleries []Gallery

  db := gg.db.Where("user_id = ?", userID)

  if err := db.Find(&galleries).Error; err != nil{
    return nil, err
  }
  return galleries, nil
}

func runGalleryValFns(gallery *Gallery, fns ...galleryValFn) error{
  for _, fn := range fns{
    err := fn(gallery)
    if err != nil{
      return err
    }
  }
  return nil
}

func (gv *galleryValidator) Create(gallery *Gallery) error{
  err := runGalleryValFns(gallery, gv.userIDRequired, gv.titleRequired)
  if err != nil{
    return err
  }
  return gv.GalleryDB.Create(gallery)
}


func (gv *galleryValidator) Update(gallery *Gallery)error{
  err := runGalleryValFns(gallery, gv.userIDRequired, gv.titleRequired)
  if err != nil{
    return err
  }
  return gv.GalleryDB.Update(gallery)
}

func (gv *galleryValidator) Delete(id uint) error{
  var gallery Gallery
  gallery.ID = id
  if err := runGalleryValFns(&gallery, gv.nonZeroID); err != nil{
    return nil
  }
  return gv.GalleryDB.Delete(gallery.ID)

}

func (gv *galleryValidator) userIDRequired(g *Gallery) error{
  if g.UserID <= 0 {
    return ErrUserIDRequired
  }
  return nil
}

func (gv *galleryValidator) titleRequired(g *Gallery)error{
  if g.Title == ""{
    return ErrTitleRequire
  }
  return nil
}

func (gv *galleryValidator) nonZeroID(gallery *Gallery)error{
  if gallery.ID <= 0{
    return ErrIDInvalid
  }
  return nil
}
