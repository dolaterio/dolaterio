package api

import "net/http"

var (
	handler http.Handler
)

func init() {
	Initialize()
	ConnectDb()
	handler, _ = Handler()
}
