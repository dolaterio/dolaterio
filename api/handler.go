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

	res.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(res)
	encoder.Encode(&jobCreatedResponse{
		Created: true,
		ID:      id,
	})
}

func tasksIndexHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	job, err := GetJob(vars["id"])

	if err != nil {
		fmt.Println(err)
	}
	// job, _ := Api.Runner.Response()

	res.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(res)
	encoder.Encode(job)
}
