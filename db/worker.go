package db

import "time"

// Worker is the model struct for workers
type Worker struct {
	ID          string            `gorethink:"id,omitempty" json:"id"`
	DockerImage string            `gorethink:"docker_image" json:"docker_image"`
	Cmd         []string          `gorethink:"cmd" json:"cmd"`
	Env         map[string]string `gorethink:"env" json:"env"`
	Timeout     time.Duration     `gorethink:"timeout,omitempty" json:"timeout"`
}

// GetWorker returns a worker from the db
func GetWorker(c *Connection, id string) (*Worker, error) {
	res, err := c.workersTable.Get(id).Run(c.s)
	defer res.Close()
	if err != nil {
		return nil, err
	}
	if res.IsNil() {
		return nil, nil
	}
	var worker Worker
	err = res.One(&worker)
	if err != nil {
		return nil, err
	}
	return &worker, nil
}

// Store inserts the worker into the db
func (worker *Worker) Store(c *Connection) error {
	res, err := c.workersTable.Insert(worker).RunWrite(c.s)
	if err != nil {
		return err
	}
	if len(res.GeneratedKeys) < 1 {
		return nil
	}
	worker.ID = res.GeneratedKeys[0]

	return nil
}

// Update updates the worker into the db.
func (worker *Worker) Update(c *Connection) error {
	_, err := c.workersTable.Update(worker).RunWrite(c.s)
	if err != nil {
		return err
	}
	return nil
}
