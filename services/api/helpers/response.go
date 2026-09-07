package helpers

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func JSONResponse(
	w http.ResponseWriter,
	status int,
	response APIResponse,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(response)
}

func Success(
	w http.ResponseWriter,
	status int,
	message string,
	data interface{},
) {
	JSONResponse(w, status, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(
	w http.ResponseWriter,
	status int,
	message string,
	errors interface{},
) {
	JSONResponse(w, status, APIResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}
