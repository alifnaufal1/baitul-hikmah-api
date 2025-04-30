package middleware

import (
	"blog-api/types"
	"blog-api/utils"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomMux struct {
	http.ServeMux
	middlewares []func(next http.Handler) http.Handler
}

var SECRET_KEY = []byte("my_scret_key")

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

func GenerateToken(userId int, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"role": role,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(SECRET_KEY)
}

func Auth(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/uploads/") {
			next.ServeHTTP(w, r)
			return
		}

		excludedRoutes := []string{"/", "/register", "/login"}
		for _, route := range excludedRoutes {
			if r.URL.Path == route {
				next.ServeHTTP(w, r)
				return
			}
		}

		authHeader := r.Header.Get("Authorization") 
		if authHeader == "" {
			utils.HandleAnyError("no token provided", w, 401)
			return
		}

		tokenString := strings.Split(authHeader, " ")
		if len(tokenString) != 2 || tokenString[0] != "Bearer" {
			utils.HandleAnyError("invalid token format", w, 401)
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString[1], claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return SECRET_KEY, nil
		})

		if err != nil || !token.Valid {
			if errors.Is(err, jwt.ErrTokenExpired) {
				utils.HandleAnyError("expired token", w, 401)
				return
			}
			utils.HandleAnyError("invalid token", w, 401)
			return
		}
		
		role, _ := claims["role"].(string)

		var userID int
		switch v := claims["user_id"].(type) {
		case float64:
			userID = int(v)
		case string:
			id, err := strconv.Atoi(v)
			if err == nil {
				userID = id
			}
		}

		ctx := context.WithValue(r.Context(), types.RoleKey, role)
		ctx = context.WithValue(ctx, types.UserKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
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

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		role, ok := r.Context().Value(types.RoleKey).(string)
		if !ok || role != "admin" {
			utils.HandleAnyError("access denied, admin only", w, 403)
			return
		}
		next.ServeHTTP(w, r)
	})
}