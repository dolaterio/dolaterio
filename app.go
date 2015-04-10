package main

import (
	"net/http"

	"github.com/dolaterio/dolaterio/api"
)

func main() {
	http.Handle("/", api.Api.Handler)
	http.ListenAndServe("localhost:8080", nil)
}
