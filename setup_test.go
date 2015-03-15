package dolaterio

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	testContainerEngine ContainerEngine
)

type fakeContainerEngine struct{}
type fakeContainer struct {
	req    *JobRequest
	stdout []byte
	stderr []byte
}

func (*fakeContainerEngine) Connect() error { return nil }
func (*fakeContainerEngine) BuildContainer(req *JobRequest) (Container, error) {
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

func assertString(t *testing.T, s1, s2 string) {
	if s1 != s2 {
		t.Errorf("Expected \"%s\", got \"%v\"", s1, s2)
	}
}
func assertStringContains(t *testing.T, s1, s2 string) {
	if strings.Contains(s1, s2) {
		t.Errorf("Expected \"%s\" to contain \"%v\"", s2, s1)
	}
}

func assertNil(t *testing.T, v interface{}) {
	if v != nil {
		t.Errorf("Expected \"%v\" to be nil", v)
	}
}

func assertNotNil(t *testing.T, v interface{}) {
	if v == nil {
		t.Errorf("Expected not to be nil")
	}
}

func assertMaxDuration(t *testing.T, d1, d2 time.Duration) {
	if d2 > d1 {
		t.Errorf("Expected %0.4fs to be shorter than %0.4fs", d2.Seconds(), d1.Seconds())
	}
}
