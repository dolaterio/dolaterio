package dolaterio

import "errors"

// Runner models a job runner
type Runner struct {
	engine    *ContainerEngine
	queue     chan *JobRequest
	responses chan *JobResponse
	stop      chan bool
	stopped   bool
}

// RunnerOptions models the data required to initialize a Runner
type RunnerOptions struct {
	Engine      *ContainerEngine
	Concurrency int
}

var (
	errRunnerStopped = errors.New("this runner is stopped")
)

// NewRunner build and initializes a runner
func NewRunner(options *RunnerOptions) (*Runner, error) {
	runner := &Runner{
		engine:    options.Engine,
		queue:     make(chan *JobRequest, options.Concurrency),
		responses: make(chan *JobResponse, options.Concurrency),
		stop:      make(chan bool, options.Concurrency),
	}
	for i := 0; i < options.Concurrency; i++ {
		go runner.run()
	}
	return runner, nil
}

// Process processes the job.
func (runner *Runner) Process(req *JobRequest) error {
	if runner.stopped {
		return errRunnerStopped
	}
	runner.queue <- req
	return nil
}

// Stop stops the job runner
func (runner *Runner) Stop() {
	runner.stopped = true
	for i := 0; i < cap(runner.queue); i++ {
		runner.stop <- true
	}
}

// Response processes the job.
func (runner *Runner) Response() *JobResponse {
	return <-runner.responses
}

func (runner *Runner) run() {
	for {
		select {
		case req := <-runner.queue:
			res := req.Execute(runner.engine)
			runner.responses <- res
		case <-runner.stop:
			return
		}
	}
}
