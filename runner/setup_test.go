package runner

import (
	"github.com/dolaterio/dolaterio/db"
	"github.com/dolaterio/dolaterio/docker"
	"github.com/dolaterio/dolaterio/queue"
)

var (
	dbConnection *db.Connection
	engine       *docker.Engine
	q            queue.Queue
)

func setup() {
	var err error
	dbConnection, err = db.NewConnection()
	if err != nil {
		panic(err)
	}
	engine, err = docker.NewEngine()
	if err != nil {
		panic(err)
	}
	engine.SkipPull = true
	q = newFakeQueue()
}

func clean() {
	dbConnection.Close()
	q.Close()
}

func logErrors(errors chan error) {
	go func() {
		for err := range errors {
			panic(err)
		}
	}()
}
