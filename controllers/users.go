package controlles

import (
  "net/http"
  "fmt"

  "lenslocked.com/views"
)

func NewUser() *Users{
  return &Users{
   NewView: views.NewView("bootstrap", "views/users/new.html"),
  }
}

type Users struct{
  NewView *views.View
}

func (u *Users) New(w http.ResponseWriter, r *http.Request){
  if err := u.NewView.Render(w, nil); err !=nil {
    panic(err)
  }
}

func(u *Users) Create(w http.ResponseWriter, r *http.Request){
  if err := r.ParseForm(); err != nil{
    panic(err)
  }
  fmt.Fprintln(w, r.PostForm["email"])
  fmt.Fprintln(w, r.PostForm["password"])
  fmt.Println(r.PostForm["email"])
  fmt.Println(r.PostForm["password"])
}
