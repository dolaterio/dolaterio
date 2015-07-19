package db

import (
	"time"

	"github.com/Sirupsen/logrus"
)

// Worker is the model struct for workers
type Worker struct {
	ID          string            `gorethink:"id,omitempty" json:"id"`
	DockerImage string            `gorethink:"docker_image" json:"docker_image"`
	Cmd         []string          `gorethink:"cmd" json:"cmd"`
	Env         map[string]string `gorethink:"env" json:"env"`
	Timeout     time.Duration     `gorethink:"timeout,omitempty" json:"timeout"`
}

var (
	workerLog = logrus.WithFields(logrus.Fields{
		"package": "db",
		"model":   "worker",
	})
)

// GetWorker returns a worker from the db
func GetWorker(c *Connection, id string) (*Worker, error) {
	logFields := logrus.Fields{"id": id}
	workerLog.WithFields(logFields).Info("Fetching worker")
	res, err := c.workersTable.Get(id).Run(c.s)
	defer res.Close()
	if err != nil {
		workerLog.WithFields(logFields).WithField("err", err).Error("Error fetching worker")
		return nil, err
	}
	if res.IsNil() {
		workerLog.WithFields(logFields).Debug("Worker not found")
		return nil, nil
	}
	var worker Worker
	err = res.One(&worker)
	if err != nil {
		workerLog.WithFields(logFields).WithField("err", err).Error("Error loading worker")
		return nil, err
	}
	return &worker, nil
}

// Store inserts the worker into the db
func (worker *Worker) Store(c *Connection) error {
	workerLog.Info("Storing worker")
	workerLog.WithField("worker", worker).Debug("Worker object")

	res, err := c.workersTable.Insert(worker).RunWrite(c.s)
	if err != nil {
		workerLog.WithField("err", err).Error("Error storing the worker")
		return err
	}
	if len(res.GeneratedKeys) < 1 {
		workerLog.Error("Worker not saved")
		return nil
	}
	worker.ID = res.GeneratedKeys[0]
	workerLog.WithField("worker", worker).Debug("Worker saved")

	return nil
}

// Update updates the worker into the db.
func (worker *Worker) Update(c *Connection) error {
	workerLog.WithField("id", worker.ID).Info("Updating worker")
	_, err := c.workersTable.Update(worker).RunWrite(c.s)
	if err != nil {
		workerLog.WithField("err", err).Error("Error updating the worker")
		return err
	}
	return nil
}
