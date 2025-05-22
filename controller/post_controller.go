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
  requestedPost, err := utils.CheckPost(r)
  if err != nil {
    utils.HandleAnyError(err.Error(), w, http.StatusBadRequest)
    return
  }

  userID, err := utils.GetRegisteredUserId(r)
  if err != nil || userID == 0 {
		utils.HandleAnyError(err.Error(), w, http.StatusUnauthorized)
		return
	}

  createdPost, err := repo.CreatePost(requestedPost, userID)
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
  strPostID := strings.Split(r.URL.Path, "/")[2]
  postID, _ := strconv.Atoi(strPostID)

  requestedPost, err := utils.CheckPost(r)
  if err != nil {
    utils.HandleAnyError(err.Error(), w, http.StatusBadRequest)
    return
  }

  updatedPost, err := repo.UpdatePost(postID, requestedPost)
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

  postData := types.DirName{
    IdFileName: postID,
    ImageName: updatedPost.Title,
    ImageType: "post",
  }

  fileName, err := UploadImageController(w, r, postData)
  if fileName == "" || err != nil {
    utils.HandleAnyError("error upload post image: " + err.Error(), w, http.StatusBadRequest)
    return
  }

  err = repo.UpdatePostImage(postID, fileName)
  if err != nil {
    utils.HandleAnyError("error saving post image -> "+ err.Error(), w, http.StatusInternalServerError)
    return 
  }

  updatedPost.PostImg = fileName

  utils.SuccessResponse(w, 200, updatedPost, "success update post")
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