package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dolaterio/dolaterio/db"
	"github.com/dolaterio/dolaterio/queue"
	"github.com/gorilla/mux"
)

type jobObjectRequest struct {
	DockerImage string            `json:"docker_image"`
	Stdin       string            `json:"stdin"`
	Timeout     int               `json:"timeout"`
	Env         map[string]string `json:"env"`
}

func (api *apiHandler) jobsCreateHandler(res http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var jobReq jobObjectRequest
	decoder.Decode(&jobReq)

	job := &db.Job{
		DockerImage: jobReq.DockerImage,
		Stdin:       jobReq.Stdin,
		Env:         jobReq.Env,
		Timeout:     time.Duration(jobReq.Timeout) * time.Millisecond,
		Status:      "pending",
	}
	err := job.Store(api.dbConnection)
	api.q.Enqueue(&queue.Message{
		JobID: job.ID,
	})
	if err != nil {
		renderError(res, err, 500)
		return
	}

	api.renderJob(res, job)
}

func (api *apiHandler) jobsIndexHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	job, err := db.GetJob(api.dbConnection, vars["id"])

	if err != nil {
		renderError(res, err, 500)
		return
	}

	if job == nil {
		renderError(res, errors.New("Job not found"), 404)
		return
	}

	api.renderJob(res, job)
}

func (api *apiHandler) renderJob(res http.ResponseWriter, job *db.Job) {
	res.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(res)
	encoder.Encode(job)
}
