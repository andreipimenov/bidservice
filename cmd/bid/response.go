package main

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse is struct for response if critical error occured (for example, no one source responsed or no one price provided)
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// WriteResponse marshal interface{} value to json and print it to ResponseWriter
func WriteResponse(w http.ResponseWriter, code int, data interface{}) {
	j, _ := json.Marshal(data)
	w.WriteHeader(code)
	w.Write(j)
}
