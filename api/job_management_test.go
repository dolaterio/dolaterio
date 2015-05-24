package main

import (
	"fmt"
	"testing"

	"github.com/dolaterio/dolaterio/db"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndFetchJob(t *testing.T) {
	setup()
	defer clean()

	worker := createWorker(t, `{"docker_image":"dolaterio/dummy-worker"}`)
	job := createJob(t, fmt.Sprintf(`{"worker_id": "%v"}`, worker.ID))
	job = fetchJob(t, job.ID)
	assert.Equal(t, db.StatusFinished, job.Status)
	assert.Empty(t, job.Syserr)
}

func TestCreateAndFetchJobWithStdin(t *testing.T) {
	setup()
	defer clean()

	worker := createWorker(t, `{"docker_image":"dolaterio/dummy-worker"}`)
	job := createJob(t, fmt.Sprintf(`{"worker_id": "%v","stdin":"hello world"}`, worker.ID))
	job = fetchJob(t, job.ID)

	assert.Empty(t, job.Syserr)
	assert.Contains(t, job.Stdout, "hello world")
}

func TestCreateAndFetchJobWithJobEnvVars(t *testing.T) {
	setup()
	defer clean()

	worker := createWorker(t, `{"docker_image":"dolaterio/dummy-worker"}`)
	job := createJob(t, fmt.Sprintf(`{"worker_id": "%v","env":{"HELLO":"world"}}`, worker.ID))
	job = fetchJob(t, job.ID)
	assert.Empty(t, job.Syserr)
	assert.Contains(t, job.Stdout, "HELLO: 'world'")
}

func TestCreateAndFetchJobWithWorkerEnvVars(t *testing.T) {
	setup()
	defer clean()

	worker := createWorker(t, `{"docker_image":"dolaterio/dummy-worker","env":{"HELLO":"world"}}`)
	job := createJob(t, fmt.Sprintf(`{"worker_id": "%v"}`, worker.ID))
	job = fetchJob(t, job.ID)
	assert.Empty(t, job.Syserr)
	assert.Contains(t, job.Stdout, "HELLO: 'world'")
}

func TestCreateAndFetchJobWithBothEnvVars(t *testing.T) {
	setup()
	defer clean()

	worker := createWorker(t, `{"docker_image":"dolaterio/dummy-worker","env":{"HELLO":"no"}}`)
	job := createJob(t, fmt.Sprintf(`{"worker_id": "%v","env":{"HELLO":"world"}}`, worker.ID))
	job = fetchJob(t, job.ID)
	assert.Empty(t, job.Syserr)
	assert.Contains(t, job.Stdout, "HELLO: 'world'")
}

func TestCreateAndFetchJobWithWorkerTimeout(t *testing.T) {
	setup()
	defer clean()

	worker := createWorker(t, `{"docker_image":"dolaterio/dummy-worker","timeout":10}`)
	job := createJob(t, fmt.Sprintf(`{"worker_id": "%v"}`, worker.ID))
	job = fetchJob(t, job.ID)
	assert.Equal(t, "timeout", job.Syserr)
}
