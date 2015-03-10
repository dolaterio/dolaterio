package runner

import (
	"testing"
	"time"
)

func TestSimpleProcess(t *testing.T) {
	begin := time.Now()
	runner, err := NewRunner(&RunnerOptions{
		Engine:      testContainerEngine,
		Concurrency: 1,
	})
	assertNil(t, err)
	err = runner.Process(&JobRequest{
		Image: "dolaterio/dummy-worker",
		Cmd:   []string{"sleep", "0.005"},
	})
	assertNil(t, err)
	_, err = runner.Response()
	assertNil(t, err)
	assertMinDuration(t, 5*time.Millisecond, time.Since(begin))
}

func TestParallelProcess(t *testing.T) {
	begin := time.Now()
	runner, err := NewRunner(&RunnerOptions{
		Engine:      testContainerEngine,
		Concurrency: 2,
	})
	assertNil(t, err)
	err = runner.Process(&JobRequest{
		Image: "dolaterio/dummy-worker",
		Cmd:   []string{"sleep", "0.005"},
	})
	err = runner.Process(&JobRequest{
		Image: "dolaterio/dummy-worker",
		Cmd:   []string{"sleep", "0.005"},
	})
	assertNil(t, err)
	_, err = runner.Response()
	assertNil(t, err)
	_, err = runner.Response()
	assertNil(t, err)
	assertMaxDuration(t, 9*time.Millisecond, time.Since(begin))
}
