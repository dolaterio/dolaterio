package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/dolaterio/dolaterio/db"
	"github.com/dolaterio/dolaterio/docker"
	"github.com/dolaterio/dolaterio/queue"
)

var (
	bind = flag.String("bind", "127.0.0.1", "Bind IP")
	port = flag.String("port", "8080", "API port")
)

func main() {
	flag.Parse()

	dbConnection, err := db.NewConnection()
	if err != nil {
		log.Fatal("Failure to connect to db: ", err)
	}

	q, err := queue.NewRedisQueue()
	if err != nil {
		log.Fatal("Failure to connect to the queue: ", err)
	}

	engine, err := docker.NewEngine(&docker.EngineConfig{})
	if err != nil {
		log.Fatal("Failure to connect to docker: ", err)
	}

	handler := &apiHandler{
		engine:       engine,
		q:            q,
		dbConnection: dbConnection,
	}

	http.Handle("/", handler.rootHandler())
	address := fmt.Sprintf("%v:%v", *bind, *port)

	fmt.Printf("Serving dolater.io api on %v\n", address)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatalf("Failure serving api on %v: ", address, err)
	}

}
