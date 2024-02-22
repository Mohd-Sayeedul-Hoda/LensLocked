package controllers

import (
  "net/http"
  "strconv"

  "lenslocked.com/context"
  "lenslocked.com/models"
  "lenslocked.com/views"


  "github.com/gorilla/mux"
)

const(
  ShowGallery = "show_gallery"
  )

type Galleries struct {
  New *views.View
  ShowView *views.View
  EditView *views.View
  gs models.GalleryService
  r *mux.Router
}


type GalleryForm struct{
  Title string `schema:"title"`
}

func NewGalleries(gs models.GalleryService, r *mux.Router) *Galleries{
  return &Galleries{
    New: views.NewView("bootstrap", "galleries/new"),
    ShowView: views.NewView("bootstrap", "galleries/show"),
    EditView: views.NewView("bootstrap", "galleries/edit"),
    gs: gs,
    r: r,
  }
}

func (g *Galleries) Create(w http.ResponseWriter, r *http.Request){
  var vd views.Data
  var form GalleryForm
  err := parseForm(r, &form);
  if err != nil{
    vd.SetAlert(err)
    g.New.Render(w, vd)
  }
  
  user := context.User(r.Context())

  gallery := models.Gallery{
    Title: form.Title,
    UserID: user.ID,
  }

  err = g.gs.Create(&gallery)
  if err != nil{
    vd.SetAlert(err)
    g.New.Render(w, vd)
    return
  }
  url, err := g.r.Get(ShowGallery).URL("id", strconv.Itoa(int(gallery.ID)))

  if err != nil{
    http.Redirect(w, r, "/", http.StatusFound)
    return 
  }

  http.Redirect(w, r, url.Path, http.StatusFound)
}

func (g *Galleries) Show(w http.ResponseWriter, r *http.Request){

  gallery, err := g.galleryByID(w, r)
  if err != nil{
    return 
  }

  var vd views.Data
  vd.Yield = gallery
  g.ShowView.Render(w, vd)
}

func (g *Galleries) galleryByID(w http.ResponseWriter, r *http.Request)(*models.Gallery, error){
  vars := mux.Vars(r)
  idStr := vars["id"]
  id, err := strconv.Atoi(idStr)
  if err != nil{
    http.Error(w, "Invalid Gallery ID", http.StatusNotFound)
    return nil, err
  }
  gallery, err := g.gs.ByID(uint(id))
  if err != nil{
    switch err{
    case models.ErrNotFound:
      http.Error(w, "Gallery not found", http.StatusNotFound)
    default:
      http.Error(w, "Whoops! Something went wrong", http.StatusInternalServerError)

  }
    return nil, err
  }
  return gallery, nil
}
