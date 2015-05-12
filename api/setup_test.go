package main

import (
	"net/http"

	"github.com/dolaterio/dolaterio/docker"
	"github.com/dolaterio/dolaterio/queue"
)

var (
	handler http.Handler
	engine  *docker.Engine
	q       queue.Queue
)

func init() {
	engine, _ = docker.NewEngine(&docker.EngineConfig{})
	q, _ = queue.NewRedisQueue()

	api := &apiHandler{
		engine: engine,
		q:      q,
	}
	handler = api.rootHandler()
}
