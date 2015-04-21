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
	err = Initialize()
	if err != nil {
		panic(err)
	}
	handler = Handler()
}
