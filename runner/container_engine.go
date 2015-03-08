package runner

type containerEngine interface {
	Run(request *JobRequest) (container, error)
}

type container interface {
	Stdout() []byte
}
