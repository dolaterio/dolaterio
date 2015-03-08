package runner

import (
	"errors"
	"strings"
)

type fakeContainerEngine struct{}
type fakeContainer struct {
	stdout []byte
}

func (*fakeContainerEngine) Connect() error { return nil }
func (*fakeContainerEngine) BuildContainer(req *JobRequest) (container, error) {
	container := &fakeContainer{}
	switch req.Cmd[0] {
	case "echo":
		container.stdout = []byte(req.Cmd[1] + "\n")
	case "env":
		container.stdout = []byte(strings.Join(req.Env.StringArray(), "\n"))
	case "cat":
		container.stdout = req.Stdin
	default:
		return nil, errors.New("Unknown command: " + req.Cmd[0])
	}

	return container, nil
}

func (container *fakeContainer) AttachStdin() error {
	return nil
}
func (container *fakeContainer) Wait() error {
	return nil
}
func (container *fakeContainer) Remove() error {
	return nil
}
func (container *fakeContainer) FetchOutput() error {
	return nil
}

func (container *fakeContainer) Stdout() []byte {
	return container.stdout
}
