package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dolaterio/dolaterio/db"
	"github.com/dolaterio/dolaterio/runner"
	"github.com/stretchr/testify/assert"
)

func createJob(t *testing.T, body string) *db.Job {
	req, _ := http.NewRequest("POST", "/v1/jobs", bytes.NewBufferString(body))

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	decoder := json.NewDecoder(w.Body)
	var job db.Job
	decoder.Decode(&job)
	return &job
}

func fetchJob(t *testing.T, id string) *db.Job {
	attempts := 50

	for attempts > 0 {
		var job db.Job

		req, err := http.NewRequest("GET", "/v1/jobs/"+id, nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		decoder := json.NewDecoder(w.Body)
		decoder.Decode(&job)

		if job.Status == db.StatusFinished {
			return &job
		}

		time.Sleep(100 * time.Millisecond)
		attempts--
	}
	t.Fatal("The job has not completed")
	return nil
}

func TestCreateAndFetchJob(t *testing.T) {
	runner := runner.NewJobRunner(&runner.JobRunnerOptions{
		Engine:      engine,
		Concurrency: 1,
		Queue:       q,
	})
	runner.Start()

	job := createJob(t, "{\"docker_image\":\"dolaterio/dummy-worker\"}")
	job = fetchJob(t, job.ID)
	assert.Equal(t, db.StatusFinished, job.Status)
}

// func TestCreateAndFetchJobWithStdin(t *testing.T) {
// 	job := createJob(t, "{\"docker_image\":\"dolaterio/dummy-worker\",\"stdin\":\"hello world\"}")
// 	job = fetchJob(t, job.ID)

// 	if !strings.Contains(job.Stdout, "hello world") {
// 		t.Fatalf("Expected %v to contain %v", job.Stdout, "hello world")
// 	}
// }

// func TestCreateAndFetchJobWithEnvVars(t *testing.T) {
// 	job := createJob(t, "{\"docker_image\":\"dolaterio/dummy-worker\",\"env\":{\"HELLO\":\"world\"}}")
// 	job = fetchJob(t, job.ID)

// 	if !strings.Contains(job.Stdout, "HELLO: 'world'") {
// 		t.Fatalf("Expected %v to contain %v", job.Stdout, "HELLO: 'world'")
// 	}
// }

// func TestCreateAndFetchJobWithTimeout(t *testing.T) {
// 	job := createJob(t, "{\"docker_image\":\"dolaterio/dummy-worker\",\"timeout\":10}")
// 	job = fetchJob(t, job.ID)

// 	if job.Syserr != "timeout" {
// 		t.Fatal("Expected job to timeout")
// 	}
// }
