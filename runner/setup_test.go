package runner

import (
	"testing"

	"github.com/dolaterio/dolaterio/db"
	"github.com/dolaterio/dolaterio/docker"
	"github.com/dolaterio/dolaterio/queue"
	"github.com/stretchr/testify/assert"
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

func logErrors(t *testing.T, errors chan error) {
	go func() {
		for err := range errors {
			assert.Nil(t, err)
		}
	}()
}
