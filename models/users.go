package models

import (
	//"errors"
	"regexp"
	"strings"

	//"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"lenslocked.com/hash"
	"lenslocked.com/rand"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

var(
  ErrNotFound modelError = "models: resource not found"
  ErrIDInvalid modelError = "models: ID provided was invalid"
  ErrPasswordIncorrect modelError = "models: incorrect password provided"
  ErrEmailRequired modelError = "models: email address is require"
  ErrEmailInvalid modelError = "models: email address is not valid"
  ErrEmailTaken modelError = "models: email address is already taken"
  ErrPasswordTooShort modelError = "models: password must be at leat 8 characters long"
  ErrPasswordRequire modelError = "models: password is require"
  ErrRememberRequire modelError = "models: remember token is required"
  ErrRememberTooShort modelError = "models: remember token must be at least 32 bytes"
)

var userPwPepper = "secret-random-string"
const hmacSecretKey = "secret-hmac-key"

type modelError string

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
  emailRegex *regexp.Regexp
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

}

func (e modelError) Error() string{
  return string(e)
}

func (e modelError) Public() string{
  s := strings.Replace(string(e), "models: ", "", 1)
  split := strings.Split(s, " ")
  tc := cases.Title(language.English)
  split[0] = tc.String(split[0])
  return strings.Join(split, " ")
}


func NewUserService(db *gorm.DB) (UserService){
  ug := &userGorm{db}
  hmac := hash.NewHmac(hmacSecretKey)
  uv := newUserValidator(ug, hmac)
  return &userService{
    UserDB: uv,
  }
}

func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator{
  return &userValidator{
    UserDB: udb,
    hmac: hmac,
    emailRegex: regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
  }
}
 
func (us *userGorm) Create(user *User) error{
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
  return us.db.Save(user).Error
}

func (us *userGorm) Delete(id uint) error {
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
    return nil, ErrPasswordIncorrect 
  default:
    return nil, err
}
}

func (uv *userValidator) ByRemember(token string) (*User, error){

  user := User{
    Remember: token,
  }
  err := runUserValFns(&user, uv.hmacRemember)
  if err != nil{
    return nil, err
  }

  return uv.UserDB.ByRemember(user.RememberHash)
}

func (uv *userValidator) Create(user *User) error{

  err := runUserValFns(user,uv.passwordRequire, uv.passwordMinLenght , uv.bcryptPassword, uv.passwordHashRequire, uv.setRememberIfUnset, uv.rememberMinBytes, uv.hmacRemember, uv.rememberHashRequire, uv.normalizeEmail, uv.requireEmail, uv.emailFormat, uv.emailIsAvail)
  if err != nil{
    return err
  }

  return uv.UserDB.Create(user)
}

func (uv *userValidator) Update(user *User) error{
  
  err := runUserValFns(user,uv.passwordMinLenght, uv.bcryptPassword, uv.passwordHashRequire, uv.rememberMinBytes, uv.hmacRemember,uv.rememberHashRequire, uv.normalizeEmail, uv.requireEmail, uv.emailFormat, uv.emailIsAvail)
  if err != nil{
    return err
  }

  return uv.UserDB.Update(user)
}

func (uv *userValidator) Delete(id uint) error {
  var user User
  user.ID = id
  err := runUserValFns(&user, uv.idGreaterThan(0))
  if err != nil{
    return err
  }
  return uv.UserDB.Delete(id)
}

func (uv *userValidator) ByEmail(email string) (*User, error){
  user := User{
    Email: email,
  }
  err := runUserValFns(&user, uv.normalizeEmail)
  if err != nil{
    return nil, err
  }
  return uv.UserDB.ByEmail(user.Email)
}

// bcryptPassword will hash a user's password with 
// an app-wide pepper bcrypt, which salts for us.

func (uv *userValidator) bcryptPassword(user *User) error{
  
  if user.Password == ""{
    return nil
  }

  pwBytes := []byte(user.Password+userPwPepper)
  hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes,bcrypt.DefaultCost)
  if err != nil{
    return err
  }
  user.PasswordHash = string(hashedBytes)
  user.Password = ""
  return nil
}

type userValFn func(*User) error

func runUserValFns(user *User, fns ...userValFn)error{
  for _, fn := range fns{
    err := fn(user)
    if err != nil{
      return err
    }
  }
  return nil
}

func (uv *userValidator) hmacRemember(user *User) error{
  if user.Remember == "" {
    return nil
  }
  user.RememberHash = uv.hmac.Hash(user.Remember)
  return nil
}

func (uv *userValidator) setRememberIfUnset(user *User) error{
  if user.Remember != ""{
    return nil
  }
  token, err := rand.RememberToken()
  if err != nil{
    return err
  }
  user.Remember = token
  return nil
}

func (uv *userValidator) idGreaterThan(n uint) userValFn{
  return userValFn(func(user *User) error{
    if user.ID <= n {
      return ErrIDInvalid 
    }
    return nil
  })
}
 
func (uv *userValidator) normalizeEmail(user *User) error{
  user.Email = strings.ToLower(user.Email)
  user.Email = strings.TrimSpace(user.Email)
  return nil
}

func (uv *userValidator) requireEmail(user *User) error{
  if user.Email == ""{
    return ErrEmailRequired
  }
  return nil
}

func (uv *userValidator) emailFormat(user *User) error{
  if user.Email == ""{
    return nil
  }
  if !uv.emailRegex.MatchString(user.Email) {
    return ErrEmailInvalid
  }
  return nil
}

func (uv *userValidator) emailIsAvail(user *User) error{
  existing, err := uv.ByEmail(user.Email)
  if err == ErrNotFound{
    return nil
  }
  // we can't continue our validaton without a 
  // successful query, so when we get error we say 
  // can't query email for internal err
  if err != nil{
    return err
  }

  if user.ID != existing.ID{
    return ErrEmailTaken
  }
  return nil
}

func (uv *userValidator) passwordMinLenght(user *User) error{
  if user.Password == "" {
    return nil
  }
  if len(user.Password) < 8{
    return ErrPasswordTooShort
  }
  return nil
}

func (uv *userValidator) passwordRequire(user *User)error{
  if user.Password == ""{
    return ErrPasswordRequire 
  }
  return nil
}

func (uv *userValidator) passwordHashRequire(user *User)error{
  if user.PasswordHash == ""{
    return ErrPasswordRequire
  }
  return nil
}

func (uv *userValidator) rememberMinBytes(user *User) error{
  if user.Remember == ""{
    return nil
  }

  n, err := rand.NBytes(user.Remember)
  if err != nil{
    return err
  }
  if n < 32 {
    return ErrRememberTooShort
  }
  return nil
}

func (uv *userValidator) rememberHashRequire(user *User) error{
  if user.RememberHash == ""{
    return ErrRememberRequire
  }
  return nil
}
