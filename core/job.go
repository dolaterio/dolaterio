package dolaterio

import "time"

// JobRequest models a request to run a job
type JobRequest struct {
	ID      string
	Image   string
	Cmd     []string
	Stdin   []byte
	Env     EnvVars
	Timeout time.Duration
}

// JobResponse models a request to run a job
type JobResponse struct {
	ID     string
	Stdout []byte
	Stderr []byte
	Error  error
}

// Execute runs the job
func (req *JobRequest) Execute(engine *ContainerEngine) *JobResponse {
	res := &JobResponse{
		ID: req.ID,
	}

	container, err := engine.BuildContainer(req)
	if err != nil {
		res.Error = err
		return res
	}
	defer container.Remove()

	done := make(chan int)
	errChn := make(chan error)
	go func() {
		err = container.AttachStdin()
		if err != nil {
			errChn <- err
		}

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
		res.Error = err
		return res

	case <-time.After(timeout):
		res.Error = errTimeout
		return res
	}

	err = container.FetchOutput()
	if err != nil {
		res.Error = err
		return res
	}

	res.Stdout = container.Stdout()
	res.Stderr = container.Stderr()
	return res
}
