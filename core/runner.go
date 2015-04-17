package dolaterio

import "errors"

// Runner models a job runner
type Runner struct {
	Responses chan *JobResponse

	engine  *ContainerEngine
	queue   chan *JobRequest
	stop    chan bool
	stopped bool
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
		Responses: make(chan *JobResponse, options.Concurrency),
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

func (runner *Runner) run() {
	for {
		select {
		case req := <-runner.queue:
			res := req.Execute(runner.engine)
			runner.Responses <- res
		case <-runner.stop:
			return
		}
	}
}
