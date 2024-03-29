package views

import (
  "bytes"
  "html/template"
  "io"
  "net/http"
  "path/filepath"
  "lenslocked.com/context"
)

var(
  LayoutDir string = "views/layouts/"
  TemplateDir string = "views/"
  TemplateExt string = ".html"
)

func layoutFiles() []string{
  files, err := filepath.Glob(LayoutDir+"*"+TemplateExt)
  if err != nil{
    panic(err)
  } 
  return files
}

func addTemplatePath(files []string){
  for i, f := range files{
    files[i] = TemplateDir + f
  }
}

func addTemplateExt(files []string){
  for i, f := range files{
    files[i] = f+TemplateExt
  }
}

func NewView(layout string, files ...string) *View{
  addTemplatePath(files)
  addTemplateExt(files)
  files = append(files, layoutFiles()...)
  t, err := template.ParseFiles(files...)
  if err != nil{
    panic(err)
  } 
  return &View{
    Template: t,
    Layout: layout,
  }
}

type View struct{
  Template *template.Template
  Layout string
}

func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {

  w.Header().Set("Content-Type", "text/html")
  var vd Data
  switch d := data.(type){
  case Data:
    vd = d
  default:
   vd = Data{
      Yield: data,
    }
  }
  // Lookup and set the sure to the User field
  vd.User = context.User(r.Context())
  var buf bytes.Buffer

  err := v.Template.ExecuteTemplate(&buf, v.Layout, vd)
  if err != nil{
    http.Error(w, "something went wrong. If the problem persists, please email supprot@lenslocked.com", http.StatusInternalServerError)
  }
  io.Copy(w, &buf)
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request){
  v.Render(w, r, nil)
}

