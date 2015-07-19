package db

import "github.com/Sirupsen/logrus"

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

var (
	jobLog = logrus.WithFields(logrus.Fields{
		"package": "db",
		"model":   "job",
	})
)

// GetJob returns a job from the db
func GetJob(c *Connection, id string) (*Job, error) {
	logFields := logrus.Fields{"id": id}
	jobLog.WithFields(logFields).Info("Fetching job")

	res, err := c.jobsTable.Get(id).Run(c.s)
	defer res.Close()
	if err != nil {
		jobLog.WithFields(logFields).WithField("err", err).Error("Error fetching job")
		return nil, err
	}
	if res.IsNil() {
		jobLog.WithFields(logFields).Debug("Job not found")
		return nil, nil
	}
	var job Job
	err = res.One(&job)
	if err != nil {
		jobLog.WithFields(logFields).WithField("err", err).Error("Error loading job")
		return nil, err
	}
	job.Worker, err = GetWorker(c, job.WorkerID)
	if err != nil {
		jobLog.WithFields(logFields).WithField("err", err).Error("Error getting job's worker")
		return nil, err
	}
	jobLog.WithFields(logFields).WithField("job", job).Debug("Retrieved job")
	return &job, nil
}

// Store inserts the job into the db
func (job *Job) Store(c *Connection) error {
	jobLog.Info("Storing job")
	jobLog.WithField("job", job).Debug("Job object")

	job.Status = StatusQueued

	res, err := c.jobsTable.Insert(job).RunWrite(c.s)
	if err != nil {
		jobLog.WithField("err", err).Error("Error storing the job")
		return err
	}
	if len(res.GeneratedKeys) < 1 {
		jobLog.Error("Job not saved")
		return nil
	}
	job.ID = res.GeneratedKeys[0]
	jobLog.WithField("job", job).Debug("Job saved")

	return nil
}

// Update updates the job into the db.
func (job *Job) Update(c *Connection) error {
	jobLog.WithField("id", job.ID).Info("Updating job")
	_, err := c.jobsTable.Update(job).RunWrite(c.s)
	if err != nil {
		jobLog.WithField("err", err).Error("Error updating the job")
		return err
	}
	return nil
}
