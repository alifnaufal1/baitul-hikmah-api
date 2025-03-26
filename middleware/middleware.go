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

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}