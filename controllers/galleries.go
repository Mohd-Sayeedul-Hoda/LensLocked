package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"lenslocked.com/context"
	"lenslocked.com/models"
	"lenslocked.com/views"
)

type Galleries struct {
  New *views.View
  ShowView *views.View
  gs models.GalleryService
}


type GalleryForm struct{
  Title string `schema:"title"`
}

func NewGalleries(gs models.GalleryService) *Galleries{
  return &Galleries{
    New: views.NewView("bootstrap", "galleries/new"),
    ShowView: views.NewView("bootstrap", "galleries/show"),
    gs: gs,
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
  fmt.Fprintln(w, gallery)
}

func (g *Galleries) Show(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)
  idStr := vars["id"]

  id, err := strconv.Atoi(idStr)
  if err != nil{
    http.Error(w, "Invalid gallery ID", http.StatusNotFound)
    return
  }

  _ = id

  gallery := models.Gallery{
    Title: "A temporary fake gallery with ID: "+ idStr, 
  }

  var vd views.Data
  vd.Yield = gallery
  g.ShowView.Render(w, vd)
}
