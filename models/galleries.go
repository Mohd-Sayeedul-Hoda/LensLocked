package models

import "github.com/jinzhu/gorm"

type Gallery struct{
  gorm.Model
  UserID uint `gorm:"not_null;index"`
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

func (gg *galleryGorm) Create(gallery *Gallery) error{
  return gg.db.Create(gallery).Error
}

func NewGalleryService(db *gorm.DB) GalleryService{
  return &galleryService{
    GalleryDB: &galleryValidator{
      GalleryDB: &galleryGorm{
	db: db,
      },
    },
  }
}

