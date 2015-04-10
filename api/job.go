package api

import "fmt"

// Job is the model struct for jobs
type Job struct {
	ID          string `gorethink:"id,omitempty" json:"id"`
	DockerImage string `gorethink:"docker_image" json:"docker_image"`
}

// CreateJob inserts the job into the db
func CreateJob(job *Job) (string, error) {
	res, err := JobTable.Insert(job).RunWrite(Session)
	if err != nil {
		return "", err
	}
	fmt.Println(res)
	if len(res.GeneratedKeys) < 1 {
		return "", nil
	}

	return res.GeneratedKeys[0], nil
}

// GetJob returns a job from the db
func GetJob(id string) (*Job, error) {
	res, err := JobTable.Get(id).Run(Session)
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
