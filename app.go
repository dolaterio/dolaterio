package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/dolaterio/dolaterio/api"
)

var (
	bind    = flag.String("bind", "127.0.0.1", "Bind IP")
	port    = flag.String("port", "8080", "API port")
	rdbHost = flag.String("rhost", "localhost", "RethinkDB host or IP")
	rdbPort = flag.String("rport", "28015", "RethinkDB port")
)

func main() {
	flag.Parse()

	err := api.Initialize()
	if err != nil {
		panic(err)
	}

	err = api.ConnectDb(*rdbHost, *rdbPort)
	if err != nil {
		panic(err)
	}

	handler, err := api.Handler()
	if err != nil {
		panic(err)
	}

	http.Handle("/", handler)
	address := fmt.Sprintf("%v:%v", *bind, *port)
	http.ListenAndServe(address, nil)
}
