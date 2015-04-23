package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/dolaterio/dolaterio/db"
)

var (
	bind = flag.String("bind", "127.0.0.1", "Bind IP")
	port = flag.String("port", "8080", "API port")
)

func main() {
	flag.Parse()

	err := db.Connect()
	if err != nil {
		log.Fatal("Failure to connect to db: ", err)
	}

	http.Handle("/", Handler())
	address := fmt.Sprintf("%v:%v", *bind, *port)

	fmt.Printf("Serving dolater.io api on %v\n", address)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatalf("Failure serving api on %v: ", address, err)
	}

}
