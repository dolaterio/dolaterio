package main

import (
	"net/http"
	"testing"

	"github.com/dolaterio/dolaterio/db"
	"github.com/dolaterio/dolaterio/docker"
	"github.com/dolaterio/dolaterio/queue"
	"github.com/dolaterio/dolaterio/runner"
	"github.com/stretchr/testify/assert"
)

var (
	handler   http.Handler
	engine    *docker.Engine
	dbConn    *db.Connection
	q         queue.Queue
	jobRunner *runner.JobRunner
)

func setup(t *testing.T) {
	engine, _ = docker.NewEngine(&docker.EngineConfig{})
	q, _ = queue.NewRedisQueue()
	dbConn, _ = db.NewConnection(&db.ConnectionConfig{})

	api := &apiHandler{
		dbConnection: dbConn,
		engine:       engine,
		q:            q,
	}
	handler = api.rootHandler()

	jobRunner = runner.NewJobRunner(&runner.JobRunnerOptions{
		Engine:       engine,
		Concurrency:  1,
		Queue:        q,
		DbConnection: dbConn,
	})
	jobRunner.Start()
	go func() {
		for err := range jobRunner.Errors {
			assert.Nil(t, err)
		}
	}()
}

func clean() {
	q.Close()
	dbConn.Close()
	jobRunner.Stop()
}
