package runner

type containerEngine interface {
	Connect() error
	BuildContainer(*JobRequest) (container, error)
}

type container interface {
	AttachStdin() error
	Wait() error
	Remove() error
	FetchOutput() error

	Stdout() []byte
}
