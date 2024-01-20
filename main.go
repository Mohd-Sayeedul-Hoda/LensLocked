package main

import (
  "fmt"
  "net/http"
  "lenslocked.com/views"

  "github.com/gorilla/mux"
)

var homeViews *views.View
var contactViews *views.View

func home(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "text/html")
  must(homeViews.Render(w, nil))
}

func contact(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "text/html")
  must(contactViews.Render(w, nil))
}

func faq(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "text/html")
  fmt.Fprint(w, "<p>Write you faq at </p> "+"<a href = \"mailto:support@lenslocked.com\">")
}


func main(){
  // I deeply beleive more you fuck around more you found out 

  homeViews = views.NewView("bootstrap", "views/home.gohtml")
  contactViews = views.NewView("bootstrap", "views/contact.gohtml")

  r := mux.NewRouter()
  r.HandleFunc("/", home)
  r.HandleFunc("/contact", contact)
  r.HandleFunc("/faq", faq)
  http.ListenAndServe(":3000", r)
}

func must(err error){
  if err != nil{
    panic(err)
  }
}
