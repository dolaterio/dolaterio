package dolaterio

import (
	"testing"
	"time"
)

func TestEcho(t *testing.T) {
	job := &Job{
		Image: "ubuntu:14.04",
		Cmd:   []string{"echo", "hello world"},
	}
	job.Run(testContainerEngine)
	assertNil(t, job.Error)
	assertString(t, "hello world\n", string(job.Stdout))
	assertString(t, "", string(job.Stderr))
}

func TestEnv(t *testing.T) {
	env := make(EnvVars, 2)
	env[0] = EnvVar{Key: "K1", Value: "V1"}
	env[1] = EnvVar{Key: "K2", Value: "V2"}
	job := &Job{
		Image: "ubuntu:14.04",
		Cmd:   []string{"env"},
		Env:   env,
	}
	job.Run(testContainerEngine)
	assertNil(t, job.Error)
	assertStringContains(t, "K1=V1\n", string(job.Stdout))
	assertStringContains(t, "K2=V2\n", string(job.Stdout))
	assertString(t, "", string(job.Stderr))
}

func TestStdin(t *testing.T) {
	job := &Job{
		Image: "ubuntu:14.04",
		Cmd:   []string{"cat"},
		Stdin: []byte("hello world\n"),
	}
	job.Run(testContainerEngine)
	assertNil(t, job.Error)
	assertString(t, "hello world\n", string(job.Stdout))
	assertString(t, "", string(job.Stderr))
}

func TestStderr(t *testing.T) {
	job := &Job{
		Image: "ubuntu:14.04",
		Cmd:   []string{"bash", "-c", "echo hello world >&2"},
	}
	job.Run(testContainerEngine)
	assertNil(t, job.Error)
	assertString(t, "", string(job.Stdout))
	assertString(t, "hello world\n", string(job.Stderr))
}

func TestTimeout(t *testing.T) {
	job := &Job{
		Image:   "ubuntu:14.04",
		Cmd:     []string{"sleep", "2000"},
		Timeout: 1 * time.Millisecond,
	}
	job.Run(testContainerEngine)
	assertString(t, errTimeout.Error(), job.Error.Error())
}
