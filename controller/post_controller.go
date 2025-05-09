package controller

import (
	"blog-api/repo"
	"blog-api/types"
	"blog-api/utils"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func PostCreateController(w http.ResponseWriter, r *http.Request) {
  title := r.FormValue("title")
  if title == "" {
    utils.HandleAnyError("title field required", w, http.StatusBadRequest)
    return
  }

  content := r.FormValue("content")
  if content == "" {
    utils.HandleAnyError("content field required", w, http.StatusBadRequest)
    return
  }

  categoryIDStr := r.FormValue("category_id")
  if categoryIDStr == "" {
    utils.HandleAnyError("category field required", w, http.StatusBadRequest)
    return
  }
  categoryID, err := strconv.Atoi(categoryIDStr)
  if err != nil {
    utils.HandleAnyError("invalid category_id", w, http.StatusBadRequest)
    return
  }

  description := r.FormValue("description")
  if description == "" {
    utils.HandleAnyError("description field required", w, http.StatusBadRequest)
    return
  }
  if len(description) > 50 {
    utils.HandleAnyError("too many description", w, http.StatusBadRequest)
    return
  }
  
  file, _, _ := r.FormFile("post-image")
  if file == nil {
    utils.HandleAnyError("file is required", w, http.StatusBadRequest)
    return 
  }
  defer file.Close()

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
    utils.HandleAnyError("error upload post image:"+ err.Error(), w, http.StatusBadRequest)
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
  strCategoryId := r.URL.Query().Get("category_id")
  categoryId, _ := strconv.Atoi(strCategoryId)
  arrayPath := strings.Split(r.URL.Path, "/")

  if len(arrayPath) == 3 {
    strID := arrayPath[2]
    postId, _ := strconv.Atoi(strID)
    post, err := repo.GetPostById(postId)
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