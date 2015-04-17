package api

import "net/http"

var (
	handler http.Handler
)

func init() {
	err := ConnectDb()
	if err != nil {
		panic(err)
	}
	Initialize()
	handler, _ = Handler()
}
