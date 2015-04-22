package dolaterio

import "time"

// Run runs the job against in the container engine
func (job *Job) Run(engine *ContainerEngine) {
	container, err := engine.BuildContainer(job)
	if err != nil {
		job.Error = err
		return
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

	timeout := job.Timeout
	if timeout == 0 {
		timeout = engine.Timeout()
	}

	select {
	case <-done:
	case err := <-errChn:
		job.Error = err
		return

	case <-time.After(timeout):
		job.Error = errTimeout
		return
	}

	err = container.FetchOutput()
	if err != nil {
		job.Error = err
	}

	job.Stdout = container.Stdout()
	job.Stderr = container.Stderr()
}
