package dolaterio

// JobRequest models a request to run a job
type JobRequest struct {
	Image string
	Cmd   []string
	Stdin []byte
	Env   EnvVars
}

// JobResponse models a request to run a job
type JobResponse struct {
	Stdout []byte
	Stderr []byte
}

// Execute runs the job
func (req *JobRequest) Execute(engine containerEngine) (*JobResponse, error) {
	container, err := engine.BuildContainer(req)
	if err != nil {
		return nil, err
	}
	defer container.Remove()

	err = container.AttachStdin()
	if err != nil {
		return nil, err
	}

	err = container.Wait()
	if err != nil {
		return nil, err
	}

	err = container.FetchOutput()
	if err != nil {
		return nil, err
	}

	return &JobResponse{
		Stdout: container.Stdout(),
		Stderr: container.Stderr(),
	}, nil
}
