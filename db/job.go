package db

import "time"

// Job is the model struct for jobs
type Job struct {
	ID          string            `gorethink:"id,omitempty" json:"id"`
	Status      string            `gorethink:"status" json:"status"`
	DockerImage string            `gorethink:"docker_image" json:"docker_image"`
	Cmd         []string          `gorethink:"cmd" json:"cmd"`
	Env         map[string]string `gorethink:"env" json:"env"`
	Stdin       string            `gorethink:"stdin" json:"stdin"`
	Stdout      string            `gorethink:"stdout" json:"stdout"`
	Stderr      string            `gorethink:"stderr" json:"stderr"`
	Timeout     time.Duration     `gorethink:"timeout,omitempty" json:"timeout"`
	Syserr      string            `gorethink:"syserr" json:"syserr"`
}

const (
	StatusQueued   = "queued"
	StatusRunning  = "running"
	StatusFinished = "finished"
)

// GetJob returns a job from the db
func GetJob(id string) (*Job, error) {
	err := Connect()
	if err != nil {
		return nil, err
	}

	res, err := jobsTable.Get(id).Run(s)
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
	return &job, nil
}

// Store inserts the job into the db
func (job *Job) Store() error {
	err := Connect()
	if err != nil {
		return err
	}
	job.Status = StatusQueued

	res, err := jobsTable.Insert(job).RunWrite(s)
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
func (job *Job) Update() error {
	err := Connect()
	if err != nil {
		return err
	}

	_, err = jobsTable.Update(job).RunWrite(s)
	if err != nil {
		return err
	}
	return nil
}
