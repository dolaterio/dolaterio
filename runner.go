package dolaterio

import "errors"

// Runner models a job runner
type Runner struct {
	engine      *ContainerEngine
	queue       chan *Job
	stop        chan bool
	concurrency int
	stopped     bool
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
		engine:      options.Engine,
		concurrency: options.Concurrency,
		queue:       make(chan *Job),
		stop:        make(chan bool),
	}
	for i := 0; i < runner.concurrency; i++ {
		go runner.run()
	}
	return runner, nil
}

// Process processes the job.
func (runner *Runner) Process(job *Job) error {
	if runner.stopped {
		return errRunnerStopped
	}
	runner.queue <- job
	return nil
}

// Stop stops the job runner
func (runner *Runner) Stop() {
	runner.stopped = true
	for i := 0; i < runner.concurrency; i++ {
		runner.stop <- true
	}
}

func (runner *Runner) run() {
	for {

		select {
		case job := <-runner.queue:
			job.Run(runner.engine)
		case <-runner.stop:
			return
		}
	}
}
