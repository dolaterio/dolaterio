package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignupAndLogin(t *testing.T) {
	Initialize()
	ConnectDb("d.lo", "28015")
	handler, _ := Handler()

	req, _ := http.NewRequest("POST", "/v1/users",
		bytes.NewBufferString("{\"username\":\"albert\",\"email\":\"admin@dolater.io\",\"password\":\"123456\"}"))

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != 201 {
		t.Errorf("User wasn't created, response code must be 201, got %v", w.Code)
	}

	req, _ = http.NewRequest("POST", "/v1/tokens",
		bytes.NewBufferString("{\"username\":\"albert\",\"password\":\"123456\"}"))

	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != 201 {
		t.Errorf("Token wasn't created, response code must be 201, got %v", w.Code)
	}

	decoder := json.NewDecoder(w.Body)
	var res struct {
		Token string `json:"token"`
	}
	decoder.Decode(&res)
	fmt.Println(res)
}
