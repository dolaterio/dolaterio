package runner

type containerEngine interface {
	Run(image string, cmd []string, env EnvVars) (container, error)
}

type container interface {
	Stdout() []byte
}
