package controllers

import (
  "fmt"
  "net/http"

  "lenslocked.com/models"
  "lenslocked.com/views"
)

type Users struct{
  NewView *views.View
  us *models.UserService
}

type SingupForm struct{
  Name string `schema:"name"`
  Email string `schema:"email"`
  Password string `schema:"password"`
}

func NewUser(us *models.UserService) *Users{
  return &Users{
   NewView: views.NewView("bootstrap", "users/new"),
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

