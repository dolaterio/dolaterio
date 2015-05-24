package db

// Job is the model struct for jobs
type Job struct {
	ID       string            `gorethink:"id,omitempty" json:"id"`
	Worker   *Worker           `gorethink:"-" json:"-"`
	WorkerID string            `gorethink:"worker_id" json:"worker_id"`
	Status   string            `gorethink:"status" json:"status"`
	Env      map[string]string `gorethink:"env" json:"env"`
	Stdin    string            `gorethink:"stdin" json:"stdin"`
	Stdout   string            `gorethink:"stdout" json:"stdout"`
	Stderr   string            `gorethink:"stderr" json:"stderr"`
	Syserr   string            `gorethink:"syserr" json:"syserr"`
}

const (
	StatusQueued   = "queued"
	StatusRunning  = "running"
	StatusFinished = "finished"
)

// GetJob returns a job from the db
func GetJob(c *Connection, id string) (*Job, error) {
	res, err := c.jobsTable.Get(id).Run(c.s)
	defer res.Close()
	if err != nil {
		return nil, err
	}
	if res.IsNil() {
		return nil, nil
	}
	var job Job
	err = res.One(&job)
	if err != nil {
		return nil, err
	}
	job.Worker, err = GetWorker(c, job.WorkerID)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// Store inserts the job into the db
func (job *Job) Store(c *Connection) error {
	job.Status = StatusQueued

	res, err := c.jobsTable.Insert(job).RunWrite(c.s)
	if err != nil {
		return err
	}
	if len(res.GeneratedKeys) < 1 {
		return nil
	}
	job.ID = res.GeneratedKeys[0]

	return nil
}

// Update updates the job into the db.
func (job *Job) Update(c *Connection) error {
	_, err := c.jobsTable.Update(job).RunWrite(c.s)
	if err != nil {
		return err
	}
	return nil
}
