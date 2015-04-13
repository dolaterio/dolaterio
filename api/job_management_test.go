package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func createJob(t *testing.T, body string) *Job {

	req, _ := http.NewRequest("POST", "/v1/jobs", bytes.NewBufferString(body))

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	decoder := json.NewDecoder(w.Body)
	var job Job
	decoder.Decode(&job)
	return &job
}

func fetchJobTimeout(t *testing.T, id string, attempts int) *Job {
	if attempts <= 0 {
		t.Fatal("The job has not completed")
	}
	req, _ := http.NewRequest("GET", "/v1/jobs/"+id, nil)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	decoder := json.NewDecoder(w.Body)
	var job Job
	decoder.Decode(&job)
	if job.Status != "completed" {
		time.Sleep(100 * time.Millisecond)
		return fetchJobTimeout(t, id, attempts-1)
	}
	return &job
}

func fetchJob(t *testing.T, id string) *Job {
	return fetchJobTimeout(t, id, 50)
}

func TestCreateAndFetchJob(t *testing.T) {
	job := createJob(t, "{\"docker_image\":\"dolaterio/dummy-worker\"}")
	job = fetchJob(t, job.ID)
	if len(job.Stdout) == 0 {
		t.Fatal("The job has no output")
	}
}

func TestCreateAndFetchJobWithStdin(t *testing.T) {
	job := createJob(t, "{\"docker_image\":\"dolaterio/dummy-worker\",\"stdin\":\"hello world\"}")
	job = fetchJob(t, job.ID)

	if !strings.Contains(job.Stdout, "hello world") {
		t.Fatalf("Expected %v to contain %v", job.Stdout, "hello world")
	}
}

func TestCreateAndFetchJobWithEnvVars(t *testing.T) {
	job := createJob(t, "{\"docker_image\":\"dolaterio/dummy-worker\",\"env\":{\"HELLO\":\"world\"}}")
	job = fetchJob(t, job.ID)

	if !strings.Contains(job.Stdout, "HELLO: 'world'") {
		t.Fatalf("Expected %v to contain %v", job.Stdout, "HELLO: 'world'")
	}
}
