package controllers

import(
  "fmt"
  "net/http"

  "lenslocked.com/models"
  "lenslocked.com/views"
)

type Galleries struct {
  New *views.View
  gs models.GalleryService
}


type GalleryForm struct{
  Title string `schema:"title"`
}

func NewGalleries(gs models.GalleryService) *Galleries{
  return &Galleries{
    New: views.NewView("bootstrap", "galleries/new"),
    gs: gs,
  }
}

func (g *Galleries) Create(w http.ResponseWriter, r *http.Request){
  // TODO: Implement this
}

