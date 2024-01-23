package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/views"
)

type Users struct{
  NewView *views.View
}

type SingupForm struct{
  Email string `schema:"email"`
  Password string `schema:"password"`
}

func NewUser() *Users{
  return &Users{
   NewView: views.NewView("bootstrap", "users/new"),
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
  fmt.Fprintln(w, "Email is", form.Email)
  fmt.Fprintln(w, "Password is", form.Password)

}

