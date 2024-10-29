package handlers

import (
	"encoding/json"
	"net/http"
)

// ApiResponse construct
type ApiResponse struct {
	Data   any `json:"data"`
	Status int `json:"status"`
	Err    any `json:"error,omitempty"`
}

// WriteJson helps to send a standard api response irrespective of the outputs.
func WriteJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Add("Content-Type", "application/json, charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ApiResponse{
		Data:   data,
		Status: status,
	})
}

// WriteError helps to send a standard error response to client.
func WriteError(w http.ResponseWriter, status int, err error) {
	w.Header().Add("Content-Type", "application/json, charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ApiResponse{
		Err:    err.Error(),
		Status: status,
	})
}
