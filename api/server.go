package api

import (
	"fmt"
	"net/http"

	"github.com/dancannon/gorethink"
	"github.com/dolaterio/dolaterio/core"
)

type apiData struct {
	Handler   http.Handler
	Engine    dolaterio.ContainerEngine
	Runner    *dolaterio.Runner
	DbSession *gorethink.Session
}

var Api = &apiData{}

func init() {
	Api.DbSession = connectDb()

	docker := &dolaterio.DockerContainerEngine{}
	err := docker.Connect()
	if err != nil {
		panic(err)
	}
	Api.Engine = docker

	runner, err := dolaterio.NewRunner(&dolaterio.RunnerOptions{
		Concurrency: 10,
		Engine:      Api.Engine,
	})
	if err != nil {
		panic(err)
	}
	Api.Runner = runner

	Api.Handler = handler()

	go func() {
		for {
			job, _ := Api.Runner.Response()
			dbJob, _ := GetJob(job.ID)
			dbJob.Stdout = string(job.Stdout)
			dbJob.Stderr = string(job.Stderr)
			SaveJob(dbJob)
			fmt.Println("Finished Job " + job.ID)
		}
	}()
}
