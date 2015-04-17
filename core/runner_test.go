package dolaterio

import (
	"testing"
	"time"
)

func TestSimpleProcess(t *testing.T) {
	runner, err := NewRunner(&RunnerOptions{
		Engine:      testContainerEngine,
		Concurrency: 10,
	})
	assertNil(t, err)
	err = runner.Process(&JobRequest{
		Image: "ubuntu:14.04",
		Cmd:   []string{"echo", "hello world"},
	})
	assertNil(t, err)
	res := <-runner.Responses
	assertNil(t, res.Error)
	assertString(t, "hello world\n", string(res.Stdout))
	runner.Stop()
}

func TestParallelProcess(t *testing.T) {
	begin := time.Now()
	runner, err := NewRunner(&RunnerOptions{
		Engine:      testContainerEngine,
		Concurrency: 10,
	})
	assertNil(t, err)
	for i := 0; i < 5; i++ {
		err = runner.Process(&JobRequest{
			Image: "ubuntu:14.04",
			Cmd:   []string{"sleep", "1"},
		})
		assertNil(t, err)
	}
	for i := 0; i < 5; i++ {
		res := <-runner.Responses
		assertNil(t, res.Error)
	}
	assertMaxDuration(t, 4*time.Second, time.Since(begin))
	runner.Stop()
}

func TestStopWhileWaiting(t *testing.T) {
	runner, err := NewRunner(&RunnerOptions{
		Engine:      testContainerEngine,
		Concurrency: 10,
	})
	assertNil(t, err)
	runner.Stop()
	err = runner.Process(&JobRequest{
		Image: "ubuntu:14.04",
		Cmd:   []string{"sleep", "1"},
	})
	assertNotNil(t, err)
}
