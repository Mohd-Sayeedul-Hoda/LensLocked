package middleware

import (
  "fmt"
  "net/http"

  "lenslocked.com/models"
  "lenslocked.com/context"
)

type RequireUser struct{
  models.UserService
}

func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc{

  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
  cookie, err := r.Cookie("remember_token")
  if err != nil{
    http.Redirect(w, r, "/login", http.StatusFound)
    return
  }
    user, err := mw.UserService.ByRemember(cookie.Value)
    if err != nil{
      http.Redirect(w, r, "/login", http.StatusFound)
      return
    }

    // using context for rembering user think so
    ctx := r.Context()
    ctx = context.WithUser(ctx, user)

    r = r.WithContext(ctx)

    fmt.Println()
    fmt.Println("User Found: ", user)
    next(w, r)
  })
}

func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc{
  return mw.ApplyFn(next.ServeHTTP)
}
