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
	mux.HandleFunc("/posts", routes.PostRoute)
	mux.HandleFunc("/categories", routes.CategoryRoute)
	
	// midlleware
	mux.RegisterMiddleware(m.Auth)
	mux.RegisterMiddleware(m.CorsMiddleware)

	server := new(http.Server)
	server.Addr = ":90"
	server.Handler = mux

	fmt.Println("starting server at localhost:90")
	server.ListenAndServe()
}