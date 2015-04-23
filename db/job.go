package db

import "time"

// Job is the model struct for jobs
type Job struct {
	ID          string            `gorethink:"id,omitempty"`
	Status      string            `gorethink:"status"`
	DockerImage string            `gorethink:"docker_image"`
	Cmd         []string          `gorethink:"cmd"`
	Env         map[string]string `gorethink:"env"`
	Stdin       string            `gorethink:"stdin"`
	Stdout      string            `gorethink:"stdout"`
	Stderr      string            `gorethink:"stderr"`
	Timeout     time.Duration     `gorethink:"timeout,omitempty"`
	Syserr      string            `gorethink:"syserr"`
}

// Store inserts the job into the db
func (job *Job) Store() error {
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

// GetJob returns a job from the db
func GetJob(id string) (*Job, error) {
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

// // UpdateJob returns a job from the db
// func UpdateJob(job *Job) (bool, error) {
// 	res, err := jobsTable.Update(job).RunWrite(s)
// 	if err != nil {
// 		return false, err
// 	}
// 	return res.Updated > 0, nil
// }
