package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dolaterio/dolaterio/db"
	"github.com/gorilla/mux"
)

type workerObjectRequest struct {
	DockerImage string            `json:"docker_image"`
	Timeout     int               `json:"timeout"`
	Env         map[string]string `json:"env"`
}

func (api *apiHandler) workersCreateHandler(res http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var workerReq workerObjectRequest
	decoder.Decode(&workerReq)

	worker := &db.Worker{
		DockerImage: workerReq.DockerImage,
		Env:         workerReq.Env,
		Timeout:     time.Duration(workerReq.Timeout) * time.Millisecond,
	}
	err := worker.Store(api.dbConnection)
	if err != nil {
		renderError(res, err, 500)
		return
	}

	res.WriteHeader(201)
	api.renderWorker(res, worker)
}

func (api *apiHandler) workersIndexHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	worker, err := db.GetWorker(api.dbConnection, vars["id"])

	if err != nil {
		renderError(res, err, 500)
		return
	}

	if worker == nil {
		renderError(res, errors.New("Worker not found"), 404)
		return
	}

	res.WriteHeader(200)
	api.renderWorker(res, worker)
}

func (api *apiHandler) renderWorker(res http.ResponseWriter, worker *db.Worker) {
	res.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(res)
	encoder.Encode(worker)
}
