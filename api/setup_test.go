package api

import "net/http"

var (
	handler http.Handler
)

func init() {
	Initialize()
	err := ConnectDb()
	if err != nil {
		panic(err)
	}
	handler, _ = Handler()
}
