package main

import (
	"encoding/json"
	"net/http"
)

// APIResponse is struct for common api response with optional code and message
type APIResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

// WriteResponse marshal interface{} value to json and print it to ResponseWriter
func WriteResponse(w http.ResponseWriter, code int, data interface{}) {
	j, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(j)
}
