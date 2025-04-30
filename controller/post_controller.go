package controller

import (
	"blog-api/repo"
	"blog-api/types"
	"blog-api/utils"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func PostCreateController(w http.ResponseWriter, r *http.Request) {
  userIDCtx := r.Context().Value(types.UserKey)
  if userIDCtx == nil {
    utils.HandleAnyError("Unauthorized", w, http.StatusUnauthorized)
    return
  }

  userID, ok := userIDCtx.(int)
  if !ok {
    utils.HandleAnyError("Invalid user ID", w, http.StatusInternalServerError)
    return
  }

  title := r.FormValue("title")
  content := r.FormValue("content")
  categoryIDStr := r.FormValue("category_id")
  description := r.FormValue("description")
  fmt.Println(len(description))
  if len(description) > 50 {
    utils.HandleAnyError("too many description", w, http.StatusBadRequest)
    return
  }

  categoryID, err := strconv.Atoi(categoryIDStr)
  if err != nil {
    utils.HandleAnyError("invalid category_id", w, http.StatusBadRequest)
    return
  }

  post := types.Post{
    Title:      title,
    Content:    content,
    CategoryID: categoryID,
    Description: description,
  }

  createdPost, err := repo.CreatePost(post, userID)
  if err != nil {
    if err == sql.ErrNoRows {
      utils.HandleDataNotFound("this category not found", w)
      return
    }
    utils.HandleAnyError("error saving post -> "+ err.Error(), w, http.StatusInternalServerError)
    return 
  }
  
  postData := types.DirName{
    IdFileName: createdPost.ID,
    ImageName: createdPost.Title,
    ImageType: "post",
  }
  
  fileName, err := UploadImageController(w, r, postData)
  if fileName == "" || err != nil {
    utils.HandleAnyError("error upload post image -> "+ err.Error(), w, http.StatusInternalServerError)
    return 
  }
  
  err = repo.UpdatePostImage(createdPost.ID, fileName)
  if err != nil {
    utils.HandleAnyError("error saving post image -> "+ err.Error(), w, http.StatusInternalServerError)
    return 
  }

  createdPost.PostImg = fileName 

  utils.SuccessResponse(w, 201, createdPost, "success create post")
}

func PostGetController(w http.ResponseWriter, r *http.Request)  {
  strId := r.URL.Query().Get("id")
  strCategoryId := r.URL.Query().Get("category_id")

  if strId != "" {
    id, _ := strconv.Atoi(strId)
    post, err := repo.GetPostById(id)
    if err != nil {
      if err == sql.ErrNoRows {
        utils.HandleDataNotFound("this post not found", w)
        return
      }
      utils.HandleAnyError("error get post -> " + err.Error(), w, http.StatusInternalServerError)
      return
    }

    utils.SuccessResponse(w, 200, post, "success fetch post")
  } else {
      categoryId, _ := strconv.Atoi(strCategoryId)
      posts, err := repo.GetAllPost(categoryId)
      if err != nil {
        if err == sql.ErrNoRows {
          utils.HandleAnyError("posts not found", w, http.StatusNotFound)
          return
        }
      utils.HandleAnyError("error get all post ->" + err.Error(), w, http.StatusInternalServerError)
      return
    }
    utils.SuccessResponse(w, 200, posts, "success fetch all post")
  }  
}

func PostUpdateController(w http.ResponseWriter, r *http.Request)  {
  var post types.Post
  if err := utils.DecodeFromRequest(r.Body, &post, w); err != nil {return}

  strId := r.URL.Query().Get("id")
  if strId == "" {
    utils.HandleAnyError("missing parameter", w, http.StatusBadRequest)
    return
  }

  id, _ := strconv.Atoi(strId)

  if post.Title == "" {
    utils.HandleAnyError("post title is required", w, http.StatusBadRequest)
    return  
  }

  if post.Content == "" {
    utils.HandleAnyError("post content is required", w, http.StatusBadRequest)
    return  
  }

  updatedPost, err := repo.UpdatePost(id, post)
  if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
      utils.HandleDataNotFound("this post not found", w)
      return
    }
    if errors.Is(err, utils.ErrCategoryNotFound) {
      utils.HandleDataNotFound("category specified not found", w)
      return
    }
    utils.HandleAnyError("error update post -> " + err.Error(), w, http.StatusInternalServerError)
    return  
  }

  utils.SuccessResponse(w, 201, updatedPost, "success update post")
}

func PostDeleteController(w http.ResponseWriter, r *http.Request)  {
  strId := r.URL.Query().Get("id")
  if strId == "" {
    utils.HandleAnyError("missing parameter", w, http.StatusBadRequest)
    return
  }

  id, _ := strconv.Atoi(strId)

  err := repo.DeletePost(id)
  if err != nil {
    if err == sql.ErrNoRows {
      utils.HandleDataNotFound("this post not found", w)
      return
    }
    utils.HandleAnyError("error delete post -> " + err.Error(), w, http.StatusInternalServerError)
    return  
  }

  utils.SuccessResponse(w, 200, nil, "success delete post")
}