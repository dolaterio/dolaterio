package runner

type fakeContainerEngine struct{}
type fakeContainer struct {
	stdout []byte
}

var (
	testContainerEngine = &fakeContainerEngine{}
)

func (*fakeContainerEngine) Run(image string, cmd []string) container {
	return &fakeContainer{}
}

func (*fakeContainer) Stdout() []byte {
	return []byte("")
}
