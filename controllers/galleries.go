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
  IndexGallery = "index_gallery"
  EditGallery = "edit_gallery"
  )

type Galleries struct {
  New *views.View
  ShowView *views.View
  EditView *views.View
  IdexView *views.View
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
    IdexView: views.NewView("bootstrap", "galleries/index"),
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
    g.New.Render(w, r, vd)
  }
  
  user := context.User(r.Context())

  gallery := models.Gallery{
    Title: form.Title,
    UserID: user.ID,
  }

  err = g.gs.Create(&gallery)
  if err != nil{
    vd.SetAlert(err)
    g.New.Render(w, r, vd)
    return
  }
  url, err := g.r.Get(EditGallery).URL("id", strconv.Itoa(int(gallery.ID)))

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
  g.ShowView.Render(w, r, vd)
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

// Get/galleriers/:id/edit
func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request){
  gallery, err := g.galleryByID(w, r)
  if err != nil{
    return 
  }

  user := context.User(r.Context())
  if gallery.UserID != user.ID {
    http.Error(w, " you do not have permission to edit " + "this gallery", http.StatusForbidden)
    return
  }
  var vd views.Data
  vd.Yield = gallery
  g.EditView.Render(w, r, vd)
}

// Post /galleries/:id/update
func(g *Galleries) Update(w http.ResponseWriter, r *http.Request){
  gallery, err := g.galleryByID(w, r)
  if err != nil{
    return 
  }

  user := context.User(r.Context())
  if gallery.UserID != user.ID{
    http.Error(w, "Gallery not found", http.StatusNotFound)
    return
  }
  
  var vd views.Data
  vd.Yield = gallery
  var form GalleryForm
  if err := parseForm(r, &form); err != nil{
    vd.SetAlert(err)
    g.EditView.Render(w, r, vd)
    return
  }
  gallery.Title = form.Title
  err = g.gs.Update(gallery)

  if err != nil{
    vd.SetAlert(err)
  }else{
    vd.Alert = &views.Alert{
      Level: views.AlertLvlSuccess,
      Message: "Gallery update successfully",
    }
  }
  
  g.EditView.Render(w, r, vd)

}

func (g *Galleries) Delete(w http.ResponseWriter, r *http.Request){

  gallery, err := g.galleryByID(w, r)
  if err != nil{
    return 
  }

  user := context.User(r.Context())
  if gallery.UserID != user.ID{
    http.Error(w, "You do not have permission to edit this gallery", http.StatusForbidden)
  }

  var vd views.Data
  err = g.gs.Delete(gallery.ID)
  if err != nil{
    vd.SetAlert(err)
    vd.Yield = gallery
    g.EditView.Render(w, r, vd)
    return
  }

  url, err := g.r.Get(IndexGallery).URL()
  if err != nil{
    http.Redirect(w, r , "/", http.StatusNotFound)
    return 
  }
  http.Redirect(w, r, url.Path, http.StatusFound)
}

func (g *Galleries) Index(w http.ResponseWriter, r *http.Request){
  user := context.User(r.Context())
  galleries, err := g.gs.ByUserID(user.ID)
  if err != nil{
    http.Error(w, "Something went wrong", http.StatusNotFound)
    return
  }
  var vd views.Data
  vd.Yield = galleries
  g.IdexView.Render(w, r, vd)
}
