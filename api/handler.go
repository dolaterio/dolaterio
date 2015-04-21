package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Handler returns the http handler to serve the dolater.io API
func Handler() http.Handler {
	r := mux.NewRouter()
	v1 := r.PathPrefix("/v1").Subrouter()

	jobs := v1.PathPrefix("/jobs").Subrouter()
	jobs.Methods("POST").HandlerFunc(jobsCreateHandler)
	jobs.Methods("GET").Path("/{id}").HandlerFunc(jobsIndexHandler)

	return r
}
