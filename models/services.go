package models

import "github.com/jinzhu/gorm"

type Services struct {
  Gallery GalleryService
  User UserService
  db *gorm.DB
}

func NewServices(connectionInfo string) (*Services, error) {
  // TODO: Implement this...
  db, err := gorm.Open("postgres", connectionInfo)
  if err != nil{
    return nil, err
  }
  db.LogMode(true)

  return &Services{
    User: NewUserService(db),
    Gallery: &galleryGorm{},
  }, nil
}
