package dolaterio

import "time"

// JobRequest models a request to run a job
type JobRequest struct {
	Image   string
	Cmd     []string
	Stdin   []byte
	Env     EnvVars
	Timeout time.Duration
}

// JobResponse models a request to run a job
type JobResponse struct {
	Stdout []byte
	Stderr []byte
}

// Execute runs the job
func (req *JobRequest) Execute(engine ContainerEngine) (*JobResponse, error) {
	container, err := engine.BuildContainer(req)
	if err != nil {
		return nil, err
	}
	defer container.Remove()

	err = container.AttachStdin()
	if err != nil {
		return nil, err
	}

	done := make(chan int)
	errChn := make(chan error)

	go func() {
		err = container.Wait()
		if err != nil {
			errChn <- err
		} else {
			done <- 1
		}
	}()

	timeout := req.Timeout
	if timeout == 0 {
		timeout = engine.Timeout()
	}

	select {
	case <-done:
	case err := <-errChn:
		return nil, err
	case <-time.After(timeout):
		return nil, errTimeout
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
