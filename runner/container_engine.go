package runner

type containerEngine interface {
	Run(image string, cmd []string) container
}

type container interface {
	Stdout() []byte
}
