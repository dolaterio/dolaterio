package api

import (
	"encoding/json"
	"net/http"

	"github.com/dolaterio/dolaterio/core"
	"github.com/gorilla/mux"
)

type jobObjectRequest struct {
	DockerImage string `json:"docker_image"`
}

type jobCreatedResponse struct {
	Created bool `json:"created"`
	ID      bool `json:"id"`
}

type jobShowResponse struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

// Handler returns the http handler to serve the dolater.io API
func Handler(runner *dolaterio.Runner) http.Handler {
	r := mux.NewRouter()
	v1 := r.PathPrefix("/v1").Subrouter()

	tasks := v1.PathPrefix("/tasks").Subrouter()
	tasks.Methods("POST").HandlerFunc(tasksCreateHandler(runner))
	tasks.Methods("GET").HandlerFunc(tasksIndexHandler(runner))
	return r
}

func tasksCreateHandler(runner *dolaterio.Runner) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var job jobObjectRequest
		decoder.Decode(&job)
		// TODO: Do something with `job`

		runner.Process(&dolaterio.JobRequest{
			Image: job.DockerImage,
		})

		rw.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(rw)
		encoder.Encode(&jobCreatedResponse{
			Created: true,
		})
	}
}
func tasksIndexHandler(runner *dolaterio.Runner) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		job, _ := runner.Response()

		rw.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(rw)
		encoder.Encode(&jobShowResponse{
			Stdout: string(job.Stdout),
			Stderr: string(job.Stderr),
		})
	}
}
