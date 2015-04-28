package runner

import (
	"testing"
	"time"

	"github.com/dolaterio/dolaterio/db"
	"github.com/dolaterio/dolaterio/docker"
	"github.com/stretchr/testify/assert"
)

func TestSimpleProcess(t *testing.T) {
	engine, err := docker.NewEngine(&docker.EngineConfig{})
	assert.Nil(t, err)

	runner, err := NewJobRunner(&JobRunnerOptions{
		Engine:      engine,
		Concurrency: 1,
	})
	assert.Nil(t, err)

	job := &db.Job{
		DockerImage: "ubuntu:14.04",
		Cmd:         []string{"echo", "hello world"},
	}

	err = runner.Process(job)
	assert.Nil(t, err)
	assert.Empty(t, string(job.Stdout))
	runner.Stop()

	assert.Equal(t, string(job.Stdout), "hello world\n")
}

func TestParallelProcess(t *testing.T) {
	begin := time.Now()

	engine, err := docker.NewEngine(&docker.EngineConfig{})
	assert.Nil(t, err)

	runner, err := NewJobRunner(&JobRunnerOptions{
		Engine:      engine,
		Concurrency: 5,
	})
	assert.Nil(t, err)

	jobs := make([]*db.Job, 5)
	for i := 0; i < 5; i++ {
		jobs[i] = &db.Job{
			DockerImage: "ubuntu:14.04",
			Cmd:         []string{"sleep", "1"},
		}
		err = runner.Process(jobs[i])
		assert.Nil(t, err)
	}
	runner.Stop() // Waits for all tasks to finish

	for i := 0; i < 5; i++ {
		assert.Empty(t, jobs[i].Syserr)
		assert.Empty(t, jobs[i].Stderr)
	}
	assert.WithinDuration(t, time.Now(), begin, 4*time.Second)
}

func TestEngineTimeoutDoesNotHang(t *testing.T) {
	begin := time.Now()

	engine, err := docker.NewEngine(&docker.EngineConfig{
		Timeout: 1 * time.Second,
	})
	assert.Nil(t, err)

	runner, err := NewJobRunner(&JobRunnerOptions{
		Engine:      engine,
		Concurrency: 1,
	})
	assert.Nil(t, err)

	job := &db.Job{
		DockerImage: "ubuntu:14.04",
		Cmd:         []string{"sleep", "10"},
	}

	err = runner.Process(job)
	assert.Nil(t, err)
	runner.Stop()

	assert.WithinDuration(t, time.Now(), begin, 4*time.Second)
}

func TestJobTimeoutDoesNotHang(t *testing.T) {
	begin := time.Now()

	engine, err := docker.NewEngine(&docker.EngineConfig{})
	assert.Nil(t, err)

	runner, err := NewJobRunner(&JobRunnerOptions{
		Engine:      engine,
		Concurrency: 1,
	})
	assert.Nil(t, err)

	job := &db.Job{
		DockerImage: "ubuntu:14.04",
		Cmd:         []string{"sleep", "10"},
		Timeout:     1 * time.Second,
	}

	err = runner.Process(job)
	assert.Nil(t, err)
	runner.Stop()

	assert.WithinDuration(t, time.Now(), begin, 4*time.Second)
}

func TestFailsProcessingAfterStop(t *testing.T) {
	engine, err := docker.NewEngine(&docker.EngineConfig{})
	assert.Nil(t, err)

	runner, err := NewJobRunner(&JobRunnerOptions{
		Engine:      engine,
		Concurrency: 1,
	})
	assert.Nil(t, err)

	runner.Stop()

	job := &db.Job{
		DockerImage: "ubuntu:14.04",
		Cmd:         []string{"echo", "hello world"},
	}

	err = runner.Process(job)
	assert.NotNil(t, err)
}
