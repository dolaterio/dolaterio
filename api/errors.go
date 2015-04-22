package main

import (
	"encoding/json"
	"net/http"
)

type errResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func renderError(res http.ResponseWriter, err error, code int) {
	errRes := &errResponse{
		Error:   true,
		Message: err.Error(),
	}
	res.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(res)
	res.WriteHeader(code)
	encoder.Encode(errRes)
}
