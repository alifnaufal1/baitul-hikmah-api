package controller

import (
	"blog-api/repo"
	"blog-api/types"
	"blog-api/utils"
	"encoding/json"
	"net/http"
)

func CreatePostController(w http.ResponseWriter, r *http.Request) {
  if r.Method != "POST" {
    utils.HandleAnyError("invalid request method", w, http.StatusBadRequest)
    return
  }
  
  var payload types.Post
  if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
    utils.HandleAnyError("error decoding request body: "+err.Error(), w, http.StatusBadRequest)
    return  
  }
  
  createdPost, err := repo.CreatePost(payload)
  if err != nil {
    utils.HandleAnyError("error saving post: "+err.Error(), w, http.StatusInternalServerError)
    return 
  }

  response := types.APIResponse {
    Code: 201,
    Results: types.Result{
      Data: createdPost,
      Message: "success create post",
    },
    Status: "success",
  }
  
  utils.JSONTemplate(w, 201, response)
}

// func GetAllPostController(w http.ResponseWriter, r *http.Request) {
//   if r.Method != "POST" {
//     utils.HandleAnyError("invalid request method", w, http.StatusBadRequest)
//     return
//   } 

//   var payload types.Post

// }