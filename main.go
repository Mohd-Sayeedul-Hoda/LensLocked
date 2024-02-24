package main

import (
  "fmt"
  "net/http"
  "os"

  "lenslocked.com/controllers"
  "lenslocked.com/views"
  "lenslocked.com/models"
  "lenslocked.com/middleware"
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

  r := mux.NewRouter()

  staticC := controllers.NewStatic()
  userC := controllers.NewUser(services.User)
  galleriesC := controllers.NewGalleries(services.Gallery, r)

  // making middleware
  requireUserMw := middleware.RequireUser{
    UserService: services.User,
  }

  newGallery := requireUserMw.Apply(galleriesC.New)
  createGallery := requireUserMw.ApplyFn(galleriesC.Create)


  r.HandleFunc("/cookietest", userC.CookieTest)

  //static routes
  r.Handle("/", staticC.Home).Methods("GET")
  r.Handle("/contact", staticC.Contact).Methods("GET")
  r.Handle("/faq", staticC.Faq).Methods("GET")

  // User routes
  r.Handle("/signup", userC.NewView).Methods("GET")
  r.HandleFunc("/signup", userC.Create).Methods("POST")
  r.Handle("/login", userC.LoginView).Methods("GET")
  r.HandleFunc("/login", userC.Login).Methods("POST")

  // Gallery routes
  r.HandleFunc("/galleries/new", newGallery).Methods("GET")
  r.HandleFunc("/galleries", createGallery).Methods("POST")
  r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Name(controllers.ShowGallery)
  r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET").
  Name(controllers.EditGallery)
  r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.ApplyFn(galleriesC.Update)).Methods("POST")
  r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.ApplyFn(galleriesC.Delete)).Methods("POST")
  r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Index)).Methods("GET").
    Name(controllers.IndexGallery)

  fmt.Println("server running on port 3000...")
  http.ListenAndServe(":3000", r)
}

func must(err error){
  if err != nil{
    panic(err)
  }
}
