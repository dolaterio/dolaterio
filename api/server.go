package api

import (
	"fmt"

	"github.com/dolaterio/dolaterio/core"
)

var (
	Runner *dolaterio.Runner
	Engine *dolaterio.ContainerEngine
)

// Initializes a new API
func Initialize() error {
	docker := &dolaterio.ContainerEngine{}
	err := docker.Connect()
	if err != nil {
		return err
	}
	Engine = docker

	runner, err := dolaterio.NewRunner(&dolaterio.RunnerOptions{
		Concurrency: 10,
		Engine:      Engine,
	})
	if err != nil {
		return err
	}
	Runner = runner

	go func() {
		for {
			job, _ := Runner.Response()
			dbJob, _ := GetJob(job.ID)
			dbJob.Stdout = string(job.Stdout)
			dbJob.Stderr = string(job.Stderr)
			SaveJob(dbJob)
			fmt.Println("Finished Job " + job.ID)
		}
	}()

	return nil
}
