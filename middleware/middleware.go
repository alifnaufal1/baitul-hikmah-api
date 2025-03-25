package middleware

import (
	"blog-api/utils"
	"net/http"
)

type CustomMux struct {
	http.ServeMux
	middlewares []func(next http.Handler) http.Handler
}

const USERNAME = "admin"
const PASSWORD = "secret123"

func (c *CustomMux) RegisterMiddleware(next func(http.Handler) http.Handler)  {
	c.middlewares = append(c.middlewares, next)
}

func (c *CustomMux) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	var current http.Handler = &c.ServeMux
	
	for _, next := range c.middlewares {
		current = next(current)
	}

	current.ServeHTTP(w, r)
}

func Auth(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    username, password, ok := r.BasicAuth()
    if !ok {
      utils.HandleAnyError("invalid authorization", w, 401)
      return
    }	
  
    isValid := (username == USERNAME) && (password == PASSWORD)
    if !isValid {
      utils.HandleAnyError("wrong username/password", w, 401)
      return  
    }
    next.ServeHTTP(w, r)
  })
}