package models

import "github.com/jinzhu/gorm"

type Gallery struct{
  gorm.Model
  UserID uint `gorm:"not_null;index"`
}

type galleryGorm struct {
  db *gorm.DB
}

type GalleryDB interface {
  Create(gallery *Gallery)error
}

type GalleryService interface {
  GalleryDB
}


func (gg *galleryGorm) Create(gallery *Gallery) error{
  // TODO: Implement this later!
  return nil
}
