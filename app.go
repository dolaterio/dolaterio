package main

import (
	"net/http"

	"github.com/dolaterio/dolaterio/api"
	"github.com/dolaterio/dolaterio/core"
)

func main() {

	docker := &dolaterio.DockerContainerEngine{}
	err := docker.Connect()
	if err != nil {
		panic(err)
	}

	runner, err := dolaterio.NewRunner(&dolaterio.RunnerOptions{
		Concurrency: 10,
		Engine:      docker,
	})
	if err != nil {
		panic(err)
	}

	http.Handle("/", api.Handler(runner))
	http.ListenAndServe("localhost:8080", nil)
}
