package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dolaterio/dolaterio/db"
	"github.com/dolaterio/dolaterio/docker"
	"github.com/dolaterio/dolaterio/queue"
	"github.com/dolaterio/dolaterio/runner"
	"github.com/stretchr/testify/assert"
)

var (
	handler   http.Handler
	engine    *docker.Engine
	dbConn    *db.Connection
	q         queue.Queue
	jobRunner *runner.JobRunner
)

func setup() {
	var err error

	engine, err = docker.NewEngine()
	if err != nil {
		panic(err)
	}
	q, err = queue.NewRedisQueue()
	if err != nil {
		panic(err)
	}
	err = q.Empty()
	if err != nil {
		panic(err)
	}
	dbConn, err = db.NewConnection()
	if err != nil {
		panic(err)
	}

	api := &apiHandler{
		dbConnection: dbConn,
		engine:       engine,
		q:            q,
	}
	handler = api.rootHandler()

	jobRunner = runner.NewJobRunner(&runner.JobRunnerOptions{
		Engine:       engine,
		Concurrency:  1,
		Queue:        q,
		DbConnection: dbConn,
	})
	jobRunner.Start()
	go func() {
		for err := range jobRunner.Errors {
			panic(err)
		}
	}()
}

func clean() {
	q.Close()
	dbConn.Close()
	jobRunner.Stop()
}

func createJob(t *testing.T, body string) (*db.Job, *httptest.ResponseRecorder) {
	req, _ := http.NewRequest("POST", "/v1/jobs", bytes.NewBufferString(body))

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	decoder := json.NewDecoder(w.Body)
	var job db.Job
	decoder.Decode(&job)
	return &job, w
}

func createWorker(t *testing.T, body string) (*db.Worker, *httptest.ResponseRecorder) {
	req, _ := http.NewRequest("POST", "/v1/workers", bytes.NewBufferString(body))

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	decoder := json.NewDecoder(w.Body)
	var worker db.Worker
	decoder.Decode(&worker)
	return &worker, w
}

func fetchJob(t *testing.T, id string) *db.Job {
	attempts := 100

	for attempts > 0 {
		var job db.Job

		req, err := http.NewRequest("GET", "/v1/jobs/"+id, nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		decoder := json.NewDecoder(w.Body)
		decoder.Decode(&job)

		if job.Status == db.StatusFinished {
			return &job
		}

		time.Sleep(100 * time.Millisecond)
		attempts--
	}
	t.Fatal("The job has not completed")
	return nil
}

func fetchWorker(t *testing.T, id string) *db.Worker {
	var worker db.Worker

	req, err := http.NewRequest("GET", "/v1/workers/"+id, nil)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	decoder := json.NewDecoder(w.Body)
	decoder.Decode(&worker)

	return &worker
}
