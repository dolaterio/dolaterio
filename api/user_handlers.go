package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type userCreateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func userCreateHandler(res http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var userReq userCreateRequest
	decoder.Decode(&userReq)

	user := &User{
		Username: userReq.Username,
		Email:    userReq.Email,
		Password: userReq.Password,
	}
	err := CreateUser(user)
	if err != nil {
		renderError(res, err, 500)
		return
	}
	res.WriteHeader(201)
	renderUser(res, user)
}

func loginUserHandler(res http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var loginReq loginRequest
	decoder.Decode(&loginReq)

	user, err := LoginUser(loginReq.Username, loginReq.Password)
	if err != nil {
		renderError(res, err, 500)
		return
	}
	if user == nil {
		renderError(res, errors.New("User not found"), 404)
		return
	}
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString([]byte(user.Username))
	if err != nil {
		renderError(res, err, 500)
		return
	}

	loginResp := &loginResponse{
		Token: tokenString,
	}

	res.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(res)
	res.WriteHeader(201)
	encoder.Encode(loginResp)
}

func renderUser(res http.ResponseWriter, user *User) {
	res.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(res)
	encoder.Encode(user)
}
