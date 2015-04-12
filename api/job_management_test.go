package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateAndFetchJob(t *testing.T) {
	Initialize()
	ConnectDb("d.lo", "28015")
	handler, _ := Handler()

	var job Job

	req, _ := http.NewRequest("POST", "/v1/tasks",
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
	time.Sleep(1500 * time.Millisecond)
	req, _ = http.NewRequest("GET", "/v1/tasks/"+job.ID, nil)

	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	decoder = json.NewDecoder(w.Body)
	decoder.Decode(&job)

	if len(job.Stdout) == 0 {
		t.Error("The job has no output")
		return
	}
}
