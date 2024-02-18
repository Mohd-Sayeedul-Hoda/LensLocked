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

func (gv *galleryValidator) Create(gallery *Gallery) error{
  err := runGalleryValFns(gallery, gv.userIDRequired, gv.titleRequired)
  if err != nil{
    return err
  }
  return gv.GalleryDB.Create(gallery)
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
