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
	dbConnection, _ = db.NewConnection(&db.ConnectionConfig{})
	engine, _ = docker.NewEngine(&docker.EngineConfig{})
	q = newFakeQueue()
}

func clean() {
	dbConnection.Close()
	q.Close()
}
