package runner

// Runner models a job runner
type Runner struct {
	Engine containerEngine
	Queue  chan *JobRequest

	Responses chan *JobResponse
}

// RunnerOptions models the data required to initialize a Runner
type RunnerOptions struct {
	Engine      containerEngine
	Concurrency int
}

// NewRunner build and initializes a runner
func NewRunner(options *RunnerOptions) (*Runner, error) {
	runner := &Runner{
		Engine:    options.Engine,
		Queue:     make(chan *JobRequest, options.Concurrency),
		Responses: make(chan *JobResponse, options.Concurrency),
	}
	for i := 0; i < options.Concurrency; i++ {
		go runner.run()
	}
	return runner, nil
}

// Process processes the job.
func (runner *Runner) Process(req *JobRequest) error {
	runner.Queue <- req
	return nil
}

// Response processes the job.
func (runner *Runner) Response() (*JobResponse, error) {
	return <-runner.Responses, nil
}

func (runner *Runner) run() {
	for {
		req := <-runner.Queue
		res, _ := req.Execute(runner.Engine)
		runner.Responses <- res
	}
}
