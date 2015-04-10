package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type jobObjectRequest struct {
	DockerImage string `json:"docker_image"`
}

type jobCreatedResponse struct {
	Created bool   `json:"created"`
	ID      string `json:"id"`
}

type jobShowResponse struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

// handler returns the http handler to serve the dolater.io API
func handler() http.Handler {
	r := mux.NewRouter()
	v1 := r.PathPrefix("/v1").Subrouter()

	tasks := v1.PathPrefix("/tasks").Subrouter()
	tasks.Methods("POST").HandlerFunc(tasksCreateHandler)
	tasks.Methods("GET").HandlerFunc(tasksIndexHandler)
	return r
}

func tasksCreateHandler(rw http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var job jobObjectRequest
	decoder.Decode(&job)
	// TODO: Do something with `job`

	id, err := CreateJob(&Job{
		DockerImage: job.DockerImage,
	})
	if err != nil {
		fmt.Println(err)
	}

	// Api.Runner.Process(&dolaterio.JobRequest{
	// 	Image: job.DockerImage,
	// })

	rw.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(rw)
	encoder.Encode(&jobCreatedResponse{
		Created: true,
		ID:      id,
	})
}

func tasksIndexHandler(rw http.ResponseWriter, r *http.Request) {
	job, _ := Api.Runner.Response()

	rw.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(rw)
	encoder.Encode(&jobShowResponse{
		Stdout: string(job.Stdout),
		Stderr: string(job.Stderr),
	})
}
