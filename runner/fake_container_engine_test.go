package runner

import "strings"

type fakeContainerEngine struct{}
type fakeContainer struct {
	stdout []byte
}

var (
	testContainerEngine = &fakeContainerEngine{}
	// testContainerEngine = &DockerContainerEngine{}
)

func (*fakeContainerEngine) Run(image string, cmd []string, env EnvVars, stdin []byte) (container, error) {
	container := &fakeContainer{}
	if cmd[0] == "echo" {
		container.stdout = []byte(cmd[1] + "\n")
	} else if cmd[0] == "env" {
		container.stdout = []byte(strings.Join(env.StringArray(), "\n"))
	} else if cmd[0] == "cat" {
		container.stdout = stdin
	}
	return container, nil
}

func (container *fakeContainer) Stdout() []byte {
	return container.stdout
}
