package main

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/dolaterio/dolaterio/core"
	"github.com/dolaterio/dolaterio/db"
	"github.com/dolaterio/dolaterio/docker"
	"github.com/dolaterio/dolaterio/queue"
)

var (
	log = logrus.WithFields(logrus.Fields{
		"package": "api",
	})
)

func main() {
	dbConnection, err := db.NewConnection()
	if err != nil {
		log.WithField("err", err).Fatal("Failure connecting to the db")
	}
	log.Debug("Connected to the db")

	q, err := queue.NewRedisQueue()
	if err != nil {
		log.WithField("err", err).Fatal("Failure connecting to the queue")
	}
	log.Debug("Connected to the queue")

	engine, err := docker.NewEngine()
	if err != nil {
		log.WithField("err", err).Fatal("Failure connecting to docker")
	}
	log.Debug("Connected to docker")

	handler := &apiHandler{
		engine:       engine,
		q:            q,
		dbConnection: dbConnection,
	}

	http.Handle("/", handler.rootHandler())
	address := fmt.Sprintf("%v:%v", core.Config.Binding, core.Config.Port)

	log.WithField("address", address).Info("Serving dolater.io api")
	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.WithField("err", err).Fatal("Failure serving the api")
	}
}
