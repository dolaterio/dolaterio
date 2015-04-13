package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dolaterio/dolaterio/core"
	"github.com/gorilla/mux"
)

type jobObjectRequest struct {
	DockerImage string `json:"docker_image"`
}

func jobsCreateHandler(res http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var jobReq jobObjectRequest
	decoder.Decode(&jobReq)

	job := &Job{
		DockerImage: jobReq.DockerImage,
	}
	err := CreateJob(job)
	if err != nil {
		renderError(res, err, 500)
		return
	}

	Runner.Process(&dolaterio.JobRequest{
		ID:    job.ID,
		Image: job.DockerImage,
	})

	renderJob(res, job)
}

func jobsIndexHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	job, err := GetJob(vars["id"])

	if err != nil {
		renderError(res, err, 500)
		return
	}

	if job == nil {
		renderError(res, errors.New("Job not found"), 404)
		return
	}

	renderJob(res, job)
}

func renderJob(res http.ResponseWriter, job *Job) {
	res.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(res)
	encoder.Encode(job)
}
