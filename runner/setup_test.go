package runner

import "os"

var (
	testContainerEngine containerEngine
)

func init() {
	if os.Getenv("USE_DOCKER") == "1" {
		testContainerEngine = &DockerContainerEngine{}
	} else {
		testContainerEngine = &fakeContainerEngine{}
	}
	err := testContainerEngine.Connect()
	if err != nil {
		panic(err)
	}
}
