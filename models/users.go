package models

import (
  "errors"
  "fmt"

  "lenslocked.com/rand"
  "lenslocked.com/hash"

  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
  "golang.org/x/crypto/bcrypt"
)

var(
  ErrNotFound = errors.New("models: resource not found")
  ErrInvalidID = errors.New("models: ID provided was invalid")
  ErrInvalidPassword = errors.New("models: incorrect password provided")
)

var userPwPepper = "secret-random-string"
const hmacSecretKey = "secret-hmac-key"

type User struct{
  gorm.Model
  Name string 
  Email string `gorm:"not null;unique_index"`
  Password string `gorm:"-"`
  PasswordHash string `gorm:"not noll"`
  Remember string `gorm:"-"`
  RememberHash string `gorm:"not null;unique_index"`
}
 
// userGorm represents our database interaction layer
// and implement the UserDB interface fully.
type userGorm struct{
  db *gorm.DB
  hmac hash.HMAC
}

// UserService is a set of methods used to manipulat
// and work with the user model
type UserService interface{
  Authenticate(email, password string) (*User, error)
  UserDB
}

type userService struct {
  UserDB
}

type userValidator struct{
  UserDB
  hmac hash.HMAC
}

// Move to top after working
type UserDB interface{
  // Method for quering for single users
  ByID(id uint) (*User, error)
  ByEmail(email string) (*User, error)
  ByRemember(token string) (*User, error)
  
  //Method for altering users
  Create(user *User) error
  Update(user *User) error
  Delete(id uint) error

  // Used to close db connection
  Close() error

  // Migrate helpers
  AutoMigrate() error
  DestructiveReset() error
}

var _ UserDB = &userGorm{}
var _ UserService = &userService{}

func NewUserService(connection string) (UserService, error){
  ug, err := newUserGorm(connection)
  if err != nil{
    return nil, err
  }
  hmac := hash.NewHmac(hmacSecretKey)
  uv := &userValidator{
    hmac: hmac,
    UserDB: ug,
  }
  return &userService{
    UserDB: uv,
  }, nil
}

 
func newUserGorm(connectinInfo string)(*userGorm, error){
  db, err := gorm.Open("postgres", connectinInfo)
  if err != nil{
    return nil, err
  } 
  db.LogMode(true)
  hmac := hash.NewHmac(hmacSecretKey)
  return &userGorm{
    db: db,
    hmac: hmac,
  }, nil
}
 
func(us *userGorm) Close() error{
  return us.db.Close()
}

func (us *userGorm) AutoMigrate() error{
  if err := us.db.AutoMigrate(&User{}).Error; err != nil{
    return err
  }
  return nil
}

func (us *userGorm) DestructiveReset() error{
  err := us.db.DropTableIfExists(&User{}).Error
  if err != nil{
    return err
  }
  return us.AutoMigrate()
}

func (us *userGorm) Create(user *User) error{
  hashedBytes, err := bcrypt.GenerateFromPassword(
    []byte(user.Password + userPwPepper), bcrypt.DefaultCost)
  if err != nil{
    return err
  }
  user.PasswordHash = string(hashedBytes)
  fmt.Println(user.PasswordHash)
  user.Password = ""

  if user.Remember == ""{
    token, err := rand.RememberToken()
    if err != nil{
      return err
    }
    user.Remember = token
  }
  // TODO: Hash the token and set it on user.RemeberHash
  user.RememberHash = us.hmac.Hash(user.Remember)
    
  return us.db.Create(user).Error
}


func first(db *gorm.DB, dst interface{})error{
  err := db.First(dst).Error
  if err == gorm.ErrRecordNotFound{
    return ErrNotFound
  }
  return err
}

func (us *userGorm) ByID(id uint) (*User, error){
  var user User
  db := us.db.Where("id = ?", id)
  err := first(db, &user)
  if err != nil{
    return nil, err
  }
  return &user, nil

}

func (us *userGorm) ByEmail(email string) (*User, error){
  var user User
  db := us.db.Where("email = ?", email)
  err := first(db, &user)
  return &user, err
}

func (us *userGorm) Update(user *User) error {
  if user.Remember != "" {
    user.RememberHash = us.hmac.Hash(user.Remember)
  }
  return us.db.Save(user).Error
}

func (us *userGorm) Delete(id uint) error {
  if id == 0 {
    return ErrInvalidID
  }
  user := User{Model: gorm.Model{ID: id}}
  return us.db.Delete(&user).Error
}

func (us *userGorm) ByRemember(rememberHash string) (*User, error){
  var user User
  err := first(us.db.Where("remember_hash = ?", rememberHash), &user)
  if err != nil{
    return nil, err
  }
  return &user, nil
}

func (us *userService) Authenticate(email, password string) (*User, error){
  foundUser, err := us.ByEmail(email)
  if err != nil{
    return nil, err
  }
  err = bcrypt.CompareHashAndPassword(
    []byte(foundUser.PasswordHash),
    []byte(password+userPwPepper))
  switch err{
  case nil:
    return foundUser, nil
  case bcrypt.ErrMismatchedHashAndPassword:
    return nil, ErrInvalidPassword
  default:
    return nil, err
}
}

func (uv *userValidator) ByRemember(token string) (*User, error){
  rememberHash := uv.hmac.Hash(token)

  return uv.UserDB.ByRemember(rememberHash)
}
