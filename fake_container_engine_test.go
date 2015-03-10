package dolaterio

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type fakeContainerEngine struct{}
type fakeContainer struct {
	req    *JobRequest
	stdout []byte
	stderr []byte
}

func (*fakeContainerEngine) Connect() error { return nil }
func (*fakeContainerEngine) BuildContainer(req *JobRequest) (container, error) {
	return &fakeContainer{req: req}, nil
}

func (*fakeContainerEngine) Timeout() time.Duration {
	return 1 * time.Second
}

func (container *fakeContainer) AttachStdin() error {
	return nil
}
func (container *fakeContainer) Wait() error {
	switch container.req.Cmd[0] {
	case "echo":
		container.stdout = []byte(container.req.Cmd[1] + "\n")
	case "bash": // So far used only to echo to stderr
		s := container.req.Cmd[2]
		container.stderr = []byte(s[5:len(s)-4] + "\n")
	case "env":
		container.stdout = []byte(strings.Join(container.req.Env.StringArray(), "\n"))
	case "cat":
		container.stdout = container.req.Stdin
	case "sleep":
		seconds, _ := strconv.ParseFloat(container.req.Cmd[1], 64)
		time.Sleep(time.Duration(seconds*1000.0) * time.Millisecond)
	default:
		return errors.New("Unknown command: " + container.req.Cmd[0])
	}
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

func (container *fakeContainer) Stderr() []byte {
	return container.stderr
}
