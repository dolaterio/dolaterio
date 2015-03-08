package main

import (
	"fmt"

	"github.com/dolaterio/dolaterio/runner"
)

func main() {
	job := runner.JobRequest{
		Image: "dolaterio/dummy-worker",
		Stdin: []byte("Hello world"),
	}
	response, err := job.Execute(&runner.DockerContainerEngine{})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(response.Stdout))
}
