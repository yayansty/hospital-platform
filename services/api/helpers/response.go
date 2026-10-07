package helpers

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Code     int         `json:"code"`
	Message  string      `json:"message,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	Errors   interface{} `json:"errors,omitempty"`
	MetaData interface{} `json:"metaData,omitempty"`
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
	metaData interface{},
) {
	JSONResponse(w, status, APIResponse{
		Code:     status,
		Message:  message,
		Data:     data,
		MetaData: metaData,
	})
}

func Error(
	w http.ResponseWriter,
	status int,
	message string,
	errors interface{},
	metaData interface{},
) {
	JSONResponse(w, status, APIResponse{
		Code:     status,
		Message:  message,
		Errors:   errors,
		MetaData: metaData,
	})
}