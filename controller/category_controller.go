package controller

import (
	"blog-api/repo"
	"blog-api/types"
	"blog-api/utils"
	"encoding/json"
	"net/http"
	"strconv"
)

func CreateCategoryController(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
    utils.HandleAnyError("invalid request method", w, http.StatusBadRequest)
    return
	}
  
  var category types.Category
  if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
    utils.HandleAnyError("error decoding request body"+err.Error(), w, http.StatusBadRequest)
    return  
  }
  
  if category.Name == "" {
    utils.HandleAnyError("category name is required", w, http.StatusBadRequest)
    return  
  }
  
  createdCategory, err := repo.CreateCategory(category)
  if err != nil {
    utils.HandleAnyError("error saving category -> "+err.Error(), w, http.StatusInternalServerError)
    return  
  }
  
  response := types.APIResponse{
    Code: 201,
    Results: types.Result{
      Message: "success create category",
      Data: createdCategory,
    },
    Status: "success",
  }
  
  utils.JSONTemplate(w, 201, response)
}

func GetAllCategoryController(w http.ResponseWriter, r *http.Request) {
  if r.Method != "GET" {
    utils.HandleAnyError("invalid request method", w, http.StatusBadRequest)
    return
  }  
  
  categories, err := repo.GetAllCategory()
  if err != nil {
    utils.HandleAnyError("error get all category -> "+err.Error(), w, http.StatusInternalServerError)
    return  
  }
  
  response := types.APIResponse{
    Code: 201,
    Results: types.Result{
      Message: "success get all category",
      Data: categories,
    },
    Status: "success",
  }
  
  utils.JSONTemplate(w, 200, response)
}
  
func UpdateCategoryController(w http.ResponseWriter, r *http.Request)  {
  if r.Method != "PUT" {
    utils.HandleAnyError("invalid request method", w, http.StatusBadRequest)
    return
  }
  
  var category types.Category
  if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
    utils.HandleAnyError("error decoding request body"+err.Error(), w, http.StatusBadRequest)
    return  
  }
  
  if category.Name == "" {
    utils.HandleAnyError("category name is required", w, http.StatusBadRequest)
    return  
  }
  
  updatedCategory, err := repo.UpdateCategory(category)
  if err != nil {
    utils.HandleAnyError("error update category -> " + err.Error(), w, http.StatusInternalServerError)
    return  
  }

  response := types.APIResponse{
    Code: 200,
    Results: types.Result{
      Message: "success update category",
      Data: updatedCategory,
    },
    Status: "success",
  }
  
  utils.JSONTemplate(w, 200, response)
}

func DeleteCategoryController(w http.ResponseWriter, r *http.Request)  {
  if r.Method != "DELETE" {
    utils.HandleAnyError("invalid request method", w, http.StatusBadRequest)
    return
  }

  id, _ := strconv.Atoi(r.URL.Query().Get("id"))
  
  if err := repo.DeleteCategory(id); err != nil {
    utils.HandleAnyError("error delete category -> " + err.Error(), w, http.StatusInternalServerError)
    return  
  }

  response := types.APIResponse{
    Code: 200,
    Results: types.Result{
      Message: "success delete category",
      Data: nil,
    },
    Status: "success",
  }
  
  utils.JSONTemplate(w, 200, response)
}