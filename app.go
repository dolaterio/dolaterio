package main

import (
	"net/http"
	"os"

	"github.com/dolaterio/dolaterio/api"
)

func main() {
	bindIp := os.Getenv("BIND_IP")
	if bindIp == "" {
		bindIp = "127.0.0.1"
	}
	http.Handle("/", api.Api.Handler)
	http.ListenAndServe(bindIp+":8080", nil)
}
