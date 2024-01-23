package main

import (
  "net/http"
  "lenslocked.com/views"
  "lenslocked.com/controllers"

  "github.com/gorilla/mux"
)

var(
  homeViews *views.View
  contactViews *views.View
)

func main(){
  // I deeply beleive more you fuck around more you found out 

  staticC := controllers.NewStatic()
  userC := controllers.NewUser()

  r := mux.NewRouter()
  r.HandleFunc("/", staticC.Home.ServeHTTP).Methods("GET")
  r.HandleFunc("/contact", staticC.Contact.ServeHTTP).Methods("GET")
  r.HandleFunc("/signup", userC.New).Methods("GET")
  r.HandleFunc("/signup", userC.Create).Methods("POST")
  http.ListenAndServe(":3000", r)
}

func must(err error){
  if err != nil{
    panic(err)
  }
}
