package main

import (
  "fmt"
  "net/http"
  "os"

  "lenslocked.com/controllers"
  "lenslocked.com/views"
  "lenslocked.com/models"
  //"lenslocked.com/rand"
  //"lenslocked.com/hash"

  "github.com/gorilla/mux"
  "github.com/joho/godotenv"
)

var _ = godotenv.Load(".env")

var (
	connectionString = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
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

  services, err := models.NewServices(connectionString)
  if err != nil{
    panic(err)
  } 

  defer services.Close()
  services.AutoMigrate()

  staticC := controllers.NewStatic()
  userC := controllers.NewUser(services.User)
  galleriesC := controllers.NewGalleries(services.Gallery)

  r := mux.NewRouter()

  r.HandleFunc("/cookietest", userC.CookieTest)

  //static routes
  r.HandleFunc("/", staticC.Home.ServeHTTP).Methods("GET")
  r.HandleFunc("/contact", staticC.Contact.ServeHTTP).Methods("GET")
  r.HandleFunc("/faq", staticC.Faq.ServeHTTP).Methods("GET")

  // User routes
  r.HandleFunc("/signup", userC.NewView.ServeHTTP).Methods("GET")
  r.HandleFunc("/signup", userC.Create).Methods("POST")
  r.HandleFunc("/login", userC.LoginView.ServeHTTP).Methods("GET")
  r.HandleFunc("/login", userC.Login).Methods("POST")

  // Gallery routes
  r.HandleFunc("/galleries/new", galleriesC.New.ServeHTTP).Methods("GET")
  r.HandleFunc("/galleries", galleriesC.Create).Methods("POST")

  fmt.Println("server running on port 3000...")
  http.ListenAndServe(":3000", r)
}

func must(err error){
  if err != nil{
    panic(err)
  }
}
