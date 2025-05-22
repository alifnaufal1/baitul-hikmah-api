package utils

import (
	"blog-api/types"
	"errors"
	"net/http"
	"strconv"
)

func CheckPost(r *http.Request) (types.Post ,error) {
	title := r.FormValue("title")
  if title == "" { return types.Post{}, errors.New("title field required") }
	
  content := r.FormValue("content")
  if content == "" { return types.Post{}, errors.New("content field required") }
	
  categoryIDStr := r.FormValue("category_id")
  if categoryIDStr == "" { return types.Post{}, errors.New("category field required") }
	categoryID, err := strconv.Atoi(categoryIDStr)
  if err != nil { return types.Post{}, errors.New("invalid category") }
	
  description := r.FormValue("description")
  if description == "" { return types.Post{}, errors.New("description field required") }
	if len(description) > 50 { return types.Post{}, errors.New("too many description") }
  
  file, _, _ := r.FormFile("post-image")
  if file == nil { return types.Post{}, errors.New("file field required") }
  defer file.Close()

	postRequested := types.Post{
		Title: title,
		Content: content,
		CategoryID: categoryID,
		Description: description,
	}

	return postRequested, nil
}