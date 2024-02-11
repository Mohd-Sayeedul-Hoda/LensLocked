package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/models"
	"lenslocked.com/rand"
	"lenslocked.com/views"
)

type Users struct{
  NewView *views.View
  LoginView *views.View
  us models.UserService
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

func NewUser(us models.UserService) *Users{
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

  err = u.signIn(w, &user)
  if err != nil{
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  http.Redirect(w, r, "/cookietest", http.StatusFound)
}

func (u *Users) Login(w http.ResponseWriter, r *http.Request){
  form := LoginForm{}
  err := parseForm(r, &form)
  if err != nil{
    panic(err)
  }
  user, err := u.us.Authenticate(form.Email, form.Password)
  if err != nil{
  switch err{
  case models.ErrNotFound:
    fmt.Fprintln(w, "Invalid email address")
  case models.ErrPasswordIncorrect:
    fmt.Fprintln(w, "Inavald password prvided")
  case nil:
    fmt.Println(w, user)
  default:
    http.Error(w, err.Error(), http.StatusInternalServerError)
}
    return 
  }
  err = u.signIn(w, user)
  if err != nil{
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
  http.Redirect(w, r, "/cookietest", http.StatusFound)
}

func (u *Users) signIn(w http.ResponseWriter, user *models.User) error{

  if user.Remember == ""{
    token, err := rand.RememberToken()
    if err != nil{
      return err
    }
    user.Remember = token
    err = u.us.Update(user)
    if err != nil{
      return err
    }
  }
  cookie := http.Cookie{
    Name: "remember_token",
    Value: user.Remember,
    HttpOnly: true,
  }
  http.SetCookie(w, &cookie)
  return nil
}

func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request){
  cookie, err := r.Cookie("remember_token")
  if err != nil{
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  user, err := u.us.ByRemember(cookie.Value)
  if err != nil{
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
  fmt.Fprintln(w, user)
}
