package utils

import (
	"blog-api/types"
	"encoding/json"
	"net/http"
)

func JSONTemplate(w http.ResponseWriter, statusCode int, response types.APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}