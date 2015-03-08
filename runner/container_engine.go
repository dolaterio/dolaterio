package runner

type containerEngine interface {
	Run(image string, cmd []string, env EnvVars, stdin []byte) (container, error)
}

type container interface {
	Stdout() []byte
}
