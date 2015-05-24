package runner

import (
	"testing"
	"time"

	"github.com/dolaterio/dolaterio/db"
	"github.com/dolaterio/dolaterio/queue"
	"github.com/stretchr/testify/assert"
)

func TestSimpleProcess(t *testing.T) {
	setup()
	defer clean()

	job := &db.Job{
		Worker: &db.Worker{
			DockerImage: "ubuntu:14.04",
			Cmd:         []string{"echo", "hello world"},
		},
	}
	err := job.Worker.Store(dbConnection)
	assert.Nil(t, err)
	job.WorkerID = job.Worker.ID
	err = job.Store(dbConnection)
	assert.Nil(t, err)

	runner := NewJobRunner(&JobRunnerOptions{
		DbConnection: dbConnection,
		Engine:       engine,
		Concurrency:  1,
		Queue:        q,
	})
	runner.Start()
	logErrors(runner.Errors)

	q.Enqueue(&queue.Message{JobID: job.ID})
	time.Sleep(10 * time.Millisecond)

	runner.Stop()

	job, err = db.GetJob(dbConnection, job.ID)
	assert.Nil(t, err)

	assert.Equal(t, string(job.Stdout), "hello world\n")
}

func TestParallelProcess(t *testing.T) {
	setup()
	defer clean()

	runner := NewJobRunner(&JobRunnerOptions{
		DbConnection: dbConnection,
		Engine:       engine,
		Concurrency:  5,
		Queue:        q,
	})

	jobs := make([]*db.Job, 5)
	for i := 0; i < 5; i++ {
		jobs[i] = &db.Job{
			Worker: &db.Worker{
				DockerImage: "ubuntu:14.04",
				Cmd:         []string{"sleep", "1"},
			},
		}
		jobs[i].Worker.Store(dbConnection)
		jobs[i].WorkerID = jobs[i].Worker.ID
		jobs[i].Store(dbConnection)
		q.Enqueue(&queue.Message{JobID: jobs[i].ID})
	}

	begin := time.Now()
	runner.Start()
	logErrors(runner.Errors)
	time.Sleep(1 * time.Second)
	runner.Stop()

	for i := 0; i < 5; i++ {
		jobs[i], _ = db.GetJob(dbConnection, jobs[i].ID)
		assert.Empty(t, jobs[i].Syserr)
		assert.Empty(t, jobs[i].Stderr)
		assert.Equal(t, jobs[i].Status, db.StatusFinished)
	}
	assert.WithinDuration(t, time.Now(), begin, 4*time.Second)
}

func TestEngineTimeout(t *testing.T) {
	setup()
	defer clean()
	engine.Timeout = 1 * time.Second

	job := &db.Job{
		Worker: &db.Worker{
			DockerImage: "ubuntu:14.04",
			Cmd:         []string{"sleep", "10"},
		},
	}
	err := job.Worker.Store(dbConnection)
	assert.Nil(t, err)
	job.WorkerID = job.Worker.ID
	err = job.Store(dbConnection)
	assert.Nil(t, err)

	runner := NewJobRunner(&JobRunnerOptions{
		DbConnection: dbConnection,
		Engine:       engine,
		Concurrency:  1,
		Queue:        q,
	})
	begin := time.Now()
	runner.Start()
	logErrors(runner.Errors)

	q.Enqueue(&queue.Message{JobID: job.ID})
	time.Sleep(100 * time.Millisecond)

	runner.Stop()

	job, err = db.GetJob(dbConnection, job.ID)
	assert.NotEmpty(t, job.Syserr)

	assert.WithinDuration(t, time.Now(), begin, 4*time.Second)
}

func TestJobTimeout(t *testing.T) {
	setup()
	defer clean()

	job := &db.Job{
		Worker: &db.Worker{
			DockerImage: "ubuntu:14.04",
			Cmd:         []string{"sleep", "10"},
			Timeout:     1 * time.Second,
		},
	}
	err := job.Worker.Store(dbConnection)
	assert.Nil(t, err)
	job.WorkerID = job.Worker.ID
	err = job.Store(dbConnection)
	assert.Nil(t, err)

	runner := NewJobRunner(&JobRunnerOptions{
		DbConnection: dbConnection,
		Engine:       engine,
		Concurrency:  1,
		Queue:        q,
	})
	begin := time.Now()
	runner.Start()
	logErrors(runner.Errors)

	q.Enqueue(&queue.Message{JobID: job.ID})
	time.Sleep(100 * time.Millisecond)

	runner.Stop()

	job, err = db.GetJob(dbConnection, job.ID)
	assert.NotEmpty(t, job.Syserr)

	assert.WithinDuration(t, time.Now(), begin, 4*time.Second)
}
