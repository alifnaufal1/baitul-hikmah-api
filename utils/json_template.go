package utils

import (
	"blog-api/types"
	"encoding/json"
	"io"
	"net/http"
)

func SuccessResponse(w http.ResponseWriter, statusCode int, o any, message string) {
	response := types.APIResponse{
		Code: statusCode,
		Results: types.Result{
			Message: message,
			Data: o,
		},
		Status: "success",
	}
	JSONTemplate(w, statusCode, response)
}

func JSONTemplate(w http.ResponseWriter, statusCode int, response types.APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func DecodeFromRequest(requestBody io.ReadCloser, o any, w http.ResponseWriter) error {
	if err := json.NewDecoder(requestBody).Decode(o); err != nil {
		HandleAnyError("error decoding request body -> " + err.Error(), w, http.StatusBadRequest)
		return err 
	}
	return nil
}