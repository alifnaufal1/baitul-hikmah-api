package controller

import (
	"blog-api/repo"
	"blog-api/types"
	"blog-api/utils"
	"database/sql"
	"net/http"
	"strconv"
)

func CategoryCreateController(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
    utils.HandleAnyError("invalid request method", w, http.StatusBadRequest)
    return
	}
  
  var category types.Category
  if err := utils.DecodeFromRequest(r.Body, &category, w); err != nil {return}
  
  if category.Name == "" {
    utils.HandleAnyError("category name is required", w, http.StatusBadRequest)
    return  
  }
  
  createdCategory, err := repo.CreateCategory(category)
  if err != nil {
    utils.HandleAnyError("error saving category -> "+ err.Error(), w, http.StatusInternalServerError)
    return  
  }
  
  utils.SuccessResponse(w, 201, createdCategory, "success create category")
}

func CategoryGetAllController(w http.ResponseWriter, r *http.Request) {
  if r.Method != "GET" {
    utils.HandleAnyError("invalid request method", w, http.StatusBadRequest)
    return
  }  
  
  categories, err := repo.GetAllCategory()
  if err != nil {
    utils.HandleAnyError("error get all category -> " + err.Error(), w, http.StatusInternalServerError)
    return  
  }
  if len(categories) == 0 {
    utils.HandleDataNotFound("categories is empty", w)
    return
  }

  utils.SuccessResponse(w, 200, categories, "success fetch categories")
}
  
func CategoryUpdateController(w http.ResponseWriter, r *http.Request)  {
  if r.Method != "PUT" {
    utils.HandleAnyError("invalid request method", w, http.StatusBadRequest)
    return
  }
  
  var category types.Category
  if err := utils.DecodeFromRequest(r.Body, &category, w); err != nil {return}

  strId := r.URL.Query().Get("id")
  if strId == "" {
    utils.HandleAnyError("missing parameter", w, http.StatusBadRequest)
    return
  }

  id, _ := strconv.Atoi(strId)
  
  if category.Name == "" {
    utils.HandleAnyError("category name is required", w, http.StatusBadRequest)
    return  
  }
  
  updatedCategory, err := repo.UpdateCategory(id, category)
  if err != nil {
    if err == sql.ErrNoRows {
      utils.HandleDataNotFound("this category not found", w)
      return
    }
    utils.HandleAnyError("error update category -> " + err.Error(), w, http.StatusInternalServerError)
    return
  }

  utils.SuccessResponse(w, 201, updatedCategory, "success update category")
}

func CategoryDeleteController(w http.ResponseWriter, r *http.Request)  {
  if r.Method != "DELETE" {
    utils.HandleAnyError("invalid request method", w, http.StatusBadRequest)
    return
  }
  
  strId := r.URL.Query().Get("id")
  if strId == "" {
    utils.HandleAnyError("missing parameter", w, http.StatusBadRequest)
    return
  }

  id, _ := strconv.Atoi(strId)
  
  if err := repo.DeleteCategory(id); err != nil {
    if err == sql.ErrNoRows {
      utils.HandleDataNotFound("this category not found", w)
      return
    }
    utils.HandleAnyError("error delete category -> " + err.Error(), w, http.StatusInternalServerError)
    return  
  }

  utils.SuccessResponse(w, 201, nil, "success delete category")
}