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

func TestCreateAndFetchJob(t *testing.T) {
	Initialize()
	ConnectDb()
	handler, _ := Handler()

	var job Job

	req, _ := http.NewRequest("POST", "/v1/jobs",
		bytes.NewBufferString("{\"docker_image\":\"dolaterio/dummy-worker\"}"))

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	decoder := json.NewDecoder(w.Body)
	decoder.Decode(&job)

	if job.ID == "" {
		t.Error("The new job didn't get an ID")
		return
	}

	// The dummy-worker takes a bit to finish
	time.Sleep(3000 * time.Millisecond)
	req, _ = http.NewRequest("GET", "/v1/jobs/"+job.ID, nil)

	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	decoder = json.NewDecoder(w.Body)
	decoder.Decode(&job)

	if len(job.Stdout) == 0 {
		t.Error("The job has no output")
		return
	}
}

func TestCreateAndFetchJobWithStdin(t *testing.T) {
	Initialize()
	ConnectDb()
	handler, _ := Handler()

	var job Job

	req, _ := http.NewRequest("POST", "/v1/jobs",
		bytes.NewBufferString("{\"docker_image\":\"dolaterio/dummy-worker\",\"stdin\":\"hello world\"}"))

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	decoder := json.NewDecoder(w.Body)
	decoder.Decode(&job)

	if job.ID == "" {
		t.Error("The new job didn't get an ID")
		return
	}

	// The dummy-worker takes a bit to finish
	time.Sleep(3000 * time.Millisecond)
	req, _ = http.NewRequest("GET", "/v1/jobs/"+job.ID, nil)

	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	decoder = json.NewDecoder(w.Body)
	decoder.Decode(&job)

	if !strings.Contains(job.Stdout, "hello world") {
		t.Error("job's output is expected to contain 'hello world' but it didn't.")
		return
	}
}
