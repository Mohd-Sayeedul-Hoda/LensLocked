package controllers

import (
  "fmt"
  "net/http"

  "lenslocked.com/models"
  "lenslocked.com/views"
)

type Users struct{
  NewView *views.View
  LoginView *views.View
  us *models.UserService
}

type SingupForm struct{
  Name string `schema:"name"`
  Email string `schema:"email"`
  Password string `schema:"password"`
}

type LoginForm struct {
  Email string `schema:"email"`
  Password string `schema:"password"`
}

func NewUser(us *models.UserService) *Users{
  return &Users{
    NewView: views.NewView("bootstrap", "users/new"),
    LoginView: views.NewView("bootstrap", "users/login"),
    us: us,
  }
}

func (u *Users) New(w http.ResponseWriter, r *http.Request){
  if err := u.NewView.Render(w, nil); err !=nil {
    panic(err)
  }
}

func(u *Users) Create(w http.ResponseWriter, r *http.Request){
  var form SingupForm

  if err := parseForm(r, &form); err != nil{
    panic(err)
  }
  user := models.User{
    Name: form.Name,
    Email: form.Email,
    Password: form.Password,
  }
  err := u.us.Create(&user)
  if err != nil{
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  fmt.Fprintln(w, "User name is ", user.Name)
  fmt.Fprintln(w, "User email is ", user.Email)
}

func (u *Users) Login(w http.ResponseWriter, r *http.Request){
  form := LoginForm{}
  err := parseForm(r, &form)
  if err != nil{
    panic(err)
  }
  fmt.Println(form.Email)
  fmt.Println(form.Password)
  user, err := u.us.Authenticate(form.Email, form.Password)
  if err != nil{
  switch err{
  case models.ErrNotFound:
    fmt.Fprintln(w, "Invalid email address")
  case models.ErrInvalidPassword:
    fmt.Fprintln(w, "Inavald password prvided")
  case nil:
    fmt.Println(w, user)
  default:
    http.Error(w, err.Error(), http.StatusInternalServerError)
}
    return 
  }
  cookie := http.Cookie{
    Name: "email",
    Value: user.Email,
  }
  http.SetCookie(w, &cookie)
  fmt.Fprintln(w, user)
}
