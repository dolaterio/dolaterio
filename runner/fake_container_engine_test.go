package runner

type fakeContainerEngine struct{}
type fakeContainer struct {
	stdout []byte
}

var (
	testContainerEngine = &fakeContainerEngine{}
)

func (*fakeContainerEngine) Run(image string, cmd []string) container {
	container := &fakeContainer{}
	if cmd[0] == "echo" {
		container.stdout = []byte(cmd[1])
	}
	return container
}

func (container *fakeContainer) Stdout() []byte {
	return container.stdout
}
