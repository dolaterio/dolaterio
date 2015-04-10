package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dolaterio/dolaterio/core"
	"github.com/gorilla/mux"
)

type jobObjectRequest struct {
	DockerImage string `json:"docker_image"`
}

// handler returns the http handler to serve the dolater.io API
func handler() http.Handler {
	r := mux.NewRouter()
	v1 := r.PathPrefix("/v1").Subrouter()

	tasks := v1.PathPrefix("/tasks").Subrouter()
	tasks.Methods("POST").HandlerFunc(tasksCreateHandler)
	tasks.Methods("GET").Path("/{id}").HandlerFunc(tasksIndexHandler)
	return r
}

func tasksCreateHandler(res http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var jobReq jobObjectRequest
	decoder.Decode(&jobReq)

	job := &Job{
		DockerImage: jobReq.DockerImage,
	}
	err := CreateJob(job)
	if err != nil {
		fmt.Println(err)
	}

	Api.Runner.Process(&dolaterio.JobRequest{
		ID:    job.ID,
		Image: job.DockerImage,
	})

	renderJob(res, job)
}

func tasksIndexHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	job, err := GetJob(vars["id"])

	if err != nil {
		fmt.Println(err)
	}

	renderJob(res, job)
}

func renderJob(res http.ResponseWriter, job *Job) {
	res.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(res)
	encoder.Encode(job)
}
