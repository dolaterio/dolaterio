package main

import "net/http"

var (
	handler http.Handler
)

func init() {
	handler = Handler()
}
