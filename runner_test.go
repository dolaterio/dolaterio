package dolaterio

import (
	"testing"
	"time"
)

func TestSimpleProcess(t *testing.T) {
	runner, err := NewRunner(&RunnerOptions{
		Engine:      testContainerEngine,
		Concurrency: 1,
	})
	assertNil(t, err)
	job := &Job{
		Image: "ubuntu:14.04",
		Cmd:   []string{"echo", "hello world"},
	}
	err = runner.Process(job)
	assertNil(t, err)
	runner.Stop() // Waits for all tasks to finish
	assertNil(t, job.Error)
	assertString(t, "hello world\n", string(job.Stdout))
}

func TestParallelProcess(t *testing.T) {
	begin := time.Now()
	runner, err := NewRunner(&RunnerOptions{
		Engine:      testContainerEngine,
		Concurrency: 5,
	})
	assertNil(t, err)
	jobs := make([]*Job, 5)
	for i := 0; i < 5; i++ {
		jobs[i] = &Job{
			Image: "ubuntu:14.04",
			Cmd:   []string{"sleep", "1"},
		}
		err = runner.Process(jobs[i])
		assertNil(t, err)
	}
	runner.Stop() // Waits for all tasks to finish

	for i := 0; i < 5; i++ {
		assertNil(t, jobs[i].Error)
	}
	assertMaxDuration(t, 4*time.Second, time.Since(begin))
}

func TestFailsProcessingAfterStop(t *testing.T) {
	runner, err := NewRunner(&RunnerOptions{
		Engine:      testContainerEngine,
		Concurrency: 1,
	})
	assertNil(t, err)
	runner.Stop()
	err = runner.Process(&Job{
		Image: "ubuntu:14.04",
		Cmd:   []string{"sleep", "1"},
	})
	assertNotNil(t, err)
}
