package main

import (
	"encoding/json"
	"net/http"
)

type errResponse struct {
	Error   bool   `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func renderError(res http.ResponseWriter, err error, code int) {
	renderErrorMessage(res, err.Error(), code)
}

func renderErrorMessage(res http.ResponseWriter, message string, code int) {
	errRes := &errResponse{
		Error:   true,
		Code:    code,
		Message: message,
	}
	log.WithField("error_response", errRes).Error("Sent error response")
	res.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(res)
	res.WriteHeader(code)
	encoder.Encode(errRes)
}
