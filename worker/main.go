package main

import (
	"github.com/dolaterio/dolaterio/db"
	"github.com/dolaterio/dolaterio/docker"
	"github.com/dolaterio/dolaterio/queue"
	"github.com/dolaterio/dolaterio/runner"
)

func main() {
	dbConnection, err := db.NewConnection()
	if err != nil {
		panic(err)
	}
	engine, err := docker.NewEngine()
	if err != nil {
		panic(err)
	}
	queue, err := queue.NewRedisQueue()
	if err != nil {
		panic(err)
	}

	runner := runner.NewJobRunner(&runner.JobRunnerOptions{
		DbConnection: dbConnection,
		Engine:       engine,
		Queue:        queue,
		Concurrency:  8,
	})
	runner.Start()
	done := make(chan bool, 1)
	select {
	case <-done:
		runner.Stop()
	}

}
