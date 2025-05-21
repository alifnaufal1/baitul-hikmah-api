package main

import (
	"blog-api/db"
	m "blog-api/middleware"
	"blog-api/routes"
	"fmt"
	"net/http"
)

func main() {
	db.ConnectDB()
	mux := new(m.CustomMux)
	
	mux.HandleFunc("/", routes.IndexRoute)
	mux.HandleFunc("/register", routes.RegisterHandler)
	mux.HandleFunc("/login", routes.LoginHandler)
	mux.HandleFunc("/logout", routes.LogoutHandler)
	mux.HandleFunc("/request", routes.ProtectedRouteHandler)
	
	mux.HandleFunc("/posts", routes.PostRoute)
	mux.HandleFunc("/posts/", routes.PostRoute)
	mux.HandleFunc("/categories", routes.CategoryPublicRoute)
	mux.HandleFunc("/users", routes.ProfileRoute)

	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	protectedCategoryHandler := m.Auth(m.AdminOnly(http.HandlerFunc(routes.CategoryProtectedRoute)))
	mux.Handle("/categories/manage", protectedCategoryHandler)

	// midlleware
	mux.RegisterMiddleware(m.CorsMiddleware)
	mux.RegisterMiddleware(m.Auth)
	
	server := new(http.Server)
	server.Addr = ":90"
	server.Handler = m.CorsMiddleware(mux)

	fmt.Println("starting server at localhost:90")
	server.ListenAndServe()
}