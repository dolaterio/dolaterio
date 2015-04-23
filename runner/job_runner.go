package runner

import (
	"errors"

	"github.com/dolaterio/dolaterio/db"
	"github.com/dolaterio/dolaterio/docker"
)

// JobRunner models a job runner
type JobRunner struct {
	engine      *docker.Engine
	queue       chan *db.Job
	stop        chan bool
	concurrency int
	stopped     bool
}

// JobRunnerOptions models the data required to initialize a JobRunner
type JobRunnerOptions struct {
	Engine      *docker.Engine
	Concurrency int
}

var (
	errJobRunnerStopped = errors.New("this runner is stopped")
)

// NewJobRunner build and initializes a runner
func NewJobRunner(options *JobRunnerOptions) (*JobRunner, error) {
	runner := &JobRunner{
		engine:      options.Engine,
		concurrency: options.Concurrency,
		queue:       make(chan *db.Job),
		stop:        make(chan bool),
	}
	for i := 0; i < runner.concurrency; i++ {
		go runner.run()
	}
	return runner, nil
}

// Process processes the job.
func (runner *JobRunner) Process(job *db.Job) error {
	if runner.stopped {
		return errJobRunnerStopped
	}
	runner.queue <- job
	return nil
}

// Stop stops the job runner
func (runner *JobRunner) Stop() {
	runner.stopped = true
	for i := 0; i < runner.concurrency; i++ {
		runner.stop <- true
	}
}

func (runner *JobRunner) run() {
	var err error

	for {
		select {
		case job := <-runner.queue:
			err = Run(job, runner.engine)
			if err != nil {
				job.Syserr = err.Error()
			}
		case <-runner.stop:
			return
		}
	}
}
