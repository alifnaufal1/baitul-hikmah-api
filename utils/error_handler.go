package utils

import (
	"blog-api/types"
	"fmt"
	"net/http"
)

func HandleAnyError(message string, w http.ResponseWriter, statusCode int) {
	var response types.APIResponse

	if statusCode == 500 {
		response = types.APIResponse{
			Code: statusCode,
			Results: types.Result{
				Error:   "error",
				Message: "server error",
				Data:    nil,
			},
			Status: "failed",
		}

		JSONTemplate(w, statusCode, response)
		fmt.Println("Internal Server Error:", message)
		return
	}

	response = types.APIResponse{
		Code: statusCode,
		Results: types.Result{
			Error:   "error",
			Message: message,
			Data:    nil,
		},
		Status: "failed",
	}

	fmt.Println("Error:", message)

	JSONTemplate(w, statusCode, response)
}