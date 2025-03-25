package routes

import (
	"blog-api/controller"
	"blog-api/utils"
	"net/http"
)

func IndexRoute(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "GET" {
		utils.HandleAnyError("invalid request method", w, http.StatusBadRequest)
    	return
	}
	w.Write([]byte("Hello World"))
}

func PostRoute(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		controller.PostGetController(w, r);
	case "POST":
		controller.PostCreateController(w, r);
	case "PUT":
		controller.PostUpdateController(w, r);
	case "DELETE":
		controller.PostDeleteController(w, r);
	}
}

func CategoryRoute(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		controller.CategoryGetAllController(w, r);
	case "POST":
		controller.CategoryCreateController(w, r);
	case "PUT":
		controller.CategoryUpdateController(w, r);
	case "DELETE":
		controller.CategoryDeleteController(w, r);
	}
}