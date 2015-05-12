package main

import (
	"net/http"

	"github.com/dolaterio/dolaterio/db"
	"github.com/dolaterio/dolaterio/docker"
	"github.com/dolaterio/dolaterio/queue"
)

var (
	handler http.Handler
	engine  *docker.Engine
	dbConn  *db.Connection
	q       queue.Queue
)

func init() {
	engine, _ = docker.NewEngine(&docker.EngineConfig{})
	q, _ = queue.NewRedisQueue()
	dbConn, _ = db.NewConnection(&db.ConnectionConfig{})

	api := &apiHandler{
		dbConnection: dbConn,
		engine:       engine,
		q:            q,
	}
	handler = api.rootHandler()
}
