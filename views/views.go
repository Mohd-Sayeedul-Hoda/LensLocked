package views

import (
	"html/template"
	"net/http"
	"path/filepath"
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

func (v *View) Render(w http.ResponseWriter, data interface{})error{
  return v.Template.ExecuteTemplate(w, v.Layout, data)
}


func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request){
  type Alert struct{
    Level string
    Message string
  }
  alert := Alert{
    Level: "success",
    Message: "Successfully rendered a dynamic alert!",
  }

  if err := v.Render(w, alert); err != nil{
    panic(err)
  }
}
