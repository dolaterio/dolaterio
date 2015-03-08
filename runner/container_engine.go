package runner

type containerEngine interface {
	Run(image string, cmd []string, env EnvVars) container
}

type container interface {
	Stdout() []byte
}
