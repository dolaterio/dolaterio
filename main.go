package main

import (
	"fmt"

	"github.com/dolaterio/dolaterio/runner"
)

func main() {
	engine := &runner.DockerContainerEngine{}
	engine.Connect()

	job := runner.JobRequest{
		Image: "dolaterio/dummy-worker",
		Stdin: []byte("Hello world"),
	}
	response, err := job.Execute(engine)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(response.Stdout))
}
