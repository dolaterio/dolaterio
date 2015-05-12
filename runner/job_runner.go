package runner

import (
	"errors"

	"github.com/dolaterio/dolaterio/db"
	"github.com/dolaterio/dolaterio/docker"
	"github.com/dolaterio/dolaterio/queue"
)

// JobRunner models a job runner
type JobRunner struct {
	engine      *docker.Engine
	queue       queue.Queue
	jobs        chan *db.Job
	stop        chan bool
	concurrency int
	stopped     bool
}

// JobRunnerOptions models the data required to initialize a JobRunner
type JobRunnerOptions struct {
	Engine      *docker.Engine
	Concurrency int
	Queue       queue.Queue
}

var (
	errJobRunnerStopped = errors.New("this runner is stopped")
)

// NewJobRunner build and initializes a runner
func NewJobRunner(options *JobRunnerOptions) *JobRunner {
	return &JobRunner{
		engine:      options.Engine,
		concurrency: options.Concurrency,
		queue:       options.Queue,
		jobs:        make(chan *db.Job),
		stop:        make(chan bool),
	}
}

func (runner *JobRunner) Start() {
	for i := 0; i < runner.concurrency; i++ {
		go runner.run()
	}

	go func() {
		var message *queue.Message
		var job *db.Job
		cont := true

		for cont {
			message, _ = runner.queue.Dequeue()
			if message != nil {
				job, _ = db.GetJob(message.JobID)
				runner.jobs <- job
			} else {
				cont = false
			}
		}
	}()
}

// Stop stops the job runner
func (runner *JobRunner) Stop() {
	runner.queue.Close()
	runner.stopped = true
	for i := 0; i < runner.concurrency; i++ {
		runner.stop <- true
	}
}

func (runner *JobRunner) run() {
	var err error
	var job *db.Job

	for {
		select {
		case job = <-runner.jobs:
			job.Status = db.StatusQueued
			job.Update()
			err = Run(job, runner.engine)
			if err != nil {
				job.Syserr = err.Error()
			}
			job.Status = db.StatusFinished
			job.Update()
		case <-runner.stop:
			return
		}
	}
}
