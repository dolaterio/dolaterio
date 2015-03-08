package runner

import "strings"

type fakeContainerEngine struct{}
type fakeContainer struct {
	stdout []byte
}

var (
	testContainerEngine = &fakeContainerEngine{}
)

func (*fakeContainerEngine) Run(image string, cmd []string, env EnvVars) (container, error) {
	container := &fakeContainer{}
	if cmd[0] == "echo" {
		container.stdout = []byte(cmd[1])
	} else if cmd[0] == "env" {
		container.stdout = []byte(strings.Join(env.StringArray(), "\n"))
	}
	return container, nil
}

func (container *fakeContainer) Stdout() []byte {
	return container.stdout
}
