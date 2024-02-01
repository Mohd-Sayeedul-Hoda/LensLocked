package main

import (
  "fmt"
  "net/http"
  "os"

  "lenslocked.com/controllers"
  "lenslocked.com/views"
  "lenslocked.com/models"

  "github.com/gorilla/mux"
  "github.com/joho/godotenv"
)

var _ = godotenv.Load(".env")

var (
	connectionString = fmt.Sprintf("host=%s port=%s	  user=%s password=%s dbname=%s sslmode=disable",
	      os.Getenv("host"),
	      os.Getenv("port"),
	      os.Getenv("user"),
	      os.Getenv("password"),
	      os.Getenv("dbname"),
		)
)

var(
  homeViews *views.View
  contactViews *views.View
)

func main(){
  // I deeply beleive more you fuck around more you found out 
  fmt.Println(connectionString)
  us, err := models.NewUserService(connectionString)
  if err != nil{
    panic(err)
  } 

  defer us.Close()
  us.AutoMigrate()

  staticC := controllers.NewStatic()
  userC := controllers.NewUser()

  r := mux.NewRouter()
  r.HandleFunc("/", staticC.Home.ServeHTTP).Methods("GET")
  r.HandleFunc("/contact", staticC.Contact.ServeHTTP).Methods("GET")
  r.HandleFunc("/faq", staticC.Faq.ServeHTTP).Methods("GET")
  r.HandleFunc("/signup", userC.New).Methods("GET")
  r.HandleFunc("/signup", userC.Create).Methods("POST")
  http.ListenAndServe(":3000", r)
}

func must(err error){
  if err != nil{
    panic(err)
  }
}
