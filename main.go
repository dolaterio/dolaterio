package dolaterio

import (
	"fmt"
)

func main() {
	engine := &DockerContainerEngine{}
	engine.Connect()

	job := JobRequest{
		Image: "dolaterio/dummy-worker",
		Stdin: []byte("Hello world"),
	}
	response, err := job.Execute(engine)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(response.Stdout))
}
