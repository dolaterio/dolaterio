package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dolaterio/dolaterio/core"
	"github.com/dolaterio/dolaterio/db"
	"github.com/dolaterio/dolaterio/docker"
	"github.com/dolaterio/dolaterio/queue"
)

func main() {
	dbConnection, err := db.NewConnection()
	if err != nil {
		log.Fatal("Failure to connect to db: ", err)
	}

	q, err := queue.NewRedisQueue()
	if err != nil {
		log.Fatal("Failure to connect to the queue: ", err)
	}

	engine, err := docker.NewEngine()
	if err != nil {
		log.Fatal("Failure to connect to docker: ", err)
	}

	handler := &apiHandler{
		engine:       engine,
		q:            q,
		dbConnection: dbConnection,
	}

	http.Handle("/", handler.rootHandler())
	address := fmt.Sprintf("%v:%v", core.Config.Binding, core.Config.Port)

	fmt.Printf("Serving dolater.io api on %v\n", address)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatalf("Failure serving api on %v: ", address, err)
	}
}
