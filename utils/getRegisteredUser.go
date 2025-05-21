package utils

import (
	"blog-api/types"
	"errors"
	"net/http"
)

func GetUserId(r *http.Request) (int, error) {
	userIDCtx := r.Context().Value(types.UserKey)
  if userIDCtx == nil { return 0, errors.New("unauthorized") }
	
  userID, ok := userIDCtx.(int)
  if !ok { return 0, errors.New("invalid user id") }

	return userID, nil
}