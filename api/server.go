package api

import (
	"fmt"

	"github.com/dolaterio/dolaterio/core"
)

var (
	// Runner is the dolater job runner
	Runner *dolaterio.Runner
)

// Initialize a new API
func Initialize() error {
	docker := &dolaterio.ContainerEngine{}
	err := docker.Connect()
	if err != nil {
		return err
	}

	runner, err := dolaterio.NewRunner(&dolaterio.RunnerOptions{
		Concurrency: 10,
		Engine:      docker,
	})
	if err != nil {
		return err
	}
	Runner = runner

	go func() {
		for {
			job := Runner.Response()
			dbJob, _ := GetJob(job.ID)
			dbJob.Stdout = string(job.Stdout)
			dbJob.Stderr = string(job.Stderr)
			dbJob.Status = "completed"
			if job.Error != nil {
				dbJob.Syserr = job.Error.Error()
			}
			UpdateJob(dbJob)
			fmt.Println("Finished Job " + job.ID)
		}
	}()

	return nil
}
