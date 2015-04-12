package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Handler returns the http handler to serve the dolater.io API
func Handler() (http.Handler, error) {
	r := mux.NewRouter()
	v1 := r.PathPrefix("/v1").Subrouter()

	tasks := v1.PathPrefix("/tasks").Subrouter()
	tasks.Methods("POST").HandlerFunc(tasksCreateHandler)
	tasks.Methods("GET").Path("/{id}").HandlerFunc(tasksIndexHandler)

	handler := Authenticate(r)
	return handler, nil
}
