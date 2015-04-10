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
	var job Job

	req, _ := http.NewRequest("POST", "/v1/tasks",
		bytes.NewBufferString("{\"docker_image\":\"dolaterio/dummy-worker\"}"))

	w := httptest.NewRecorder()
	Api.Handler.ServeHTTP(w, req)

	decoder := json.NewDecoder(w.Body)
	decoder.Decode(&job)

	if job.ID == "" {
		t.Error("The new job didn't get an ID")
		return
	}

	// The dummy-worker takes around 6 seconds to finish
	time.Sleep(6100 * time.Millisecond)
	req, _ = http.NewRequest("GET", "/v1/tasks/"+job.ID, nil)

	w = httptest.NewRecorder()
	Api.Handler.ServeHTTP(w, req)
	decoder = json.NewDecoder(w.Body)
	decoder.Decode(&job)

	if len(job.Stdout) == 0 {
		t.Error("The job has no output")
		return
	}
}
