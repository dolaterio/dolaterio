package runner

// JobRequest models a request to run a job
type JobRequest struct {
	Image string
	Cmd   []string
	Env   EnvVars
}

// JobResponse models a request to run a job
type JobResponse struct {
	Stdout []byte
}

// Execute runs the job
func (req *JobRequest) Execute(engine containerEngine) *JobResponse {
	container := engine.Run(req.Image, req.Cmd, req.Env)

	return &JobResponse{
		Stdout: container.Stdout(),
	}
}
