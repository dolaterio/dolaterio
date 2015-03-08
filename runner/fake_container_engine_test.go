package runner

type fakeContainerEngine struct{}
type fakeContainer struct {
	stdout []byte
}

var (
	testContainerEngine = &fakeContainerEngine{}
)

func (*fakeContainerEngine) Run(image string, cmd []string, env EnvVars) container {
	container := &fakeContainer{}
	if cmd[0] == "echo" {
		container.stdout = []byte(cmd[1])
	} else if cmd[0] == "env" {
		container.stdout = []byte(env.String())
	}
	return container
}

func (container *fakeContainer) Stdout() []byte {
	return container.stdout
}
