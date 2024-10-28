package server

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	Data   any `json:"data"`
	Status int `json:"status"`
	Err    any `json:"error,omitempty"`
}

func WriteJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Add("Content-Type", "application/json, charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ApiResponse{
		Data:   data,
		Status: status,
	})
}

func WriteError(w http.ResponseWriter, status int, err error) {
	w.Header().Add("Content-Type", "application/json, charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ApiResponse{
		Err:    err.Error(),
		Status: status,
	})
}
