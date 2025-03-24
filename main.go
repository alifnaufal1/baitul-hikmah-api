package main

import (
	"blog-api/controller"
	"blog-api/db"
	m "blog-api/middleware"
	"fmt"
	"net/http"
)

func main() {
	db.ConnectDB()
	mux := new(m.CustomMux)
	
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		w.Write([]byte("Hello World"))
	})

	// post
	mux.HandleFunc("/post", controller.CreatePostController)
	
	// category
	mux.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			controller.GetAllCategoryController(w, r);
		case "POST":
			controller.CreateCategoryController(w, r);
		case "PUT":
			controller.UpdateCategoryController(w, r);
		case "DELETE":
			controller.DeleteCategoryController(w, r);
		}
	})

	// midlleware
	mux.RegisterMiddleware(m.Auth)

	server := new(http.Server)
	server.Addr = ":90"
	server.Handler = mux

	fmt.Println("starting server at localhost:90")
	server.ListenAndServe()
}