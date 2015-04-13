package dolaterio

import (
	"testing"
	"time"
)

func TestEcho(t *testing.T) {
	req := &JobRequest{
		Image: "ubuntu:14.04",
		Cmd:   []string{"echo", "hello world"},
	}
	res := req.Execute(testContainerEngine)
	assertNil(t, res.Error)
	assertString(t, "hello world\n", string(res.Stdout))
	assertString(t, "", string(res.Stderr))
}

func TestEnv(t *testing.T) {
	env := make(EnvVars, 2)
	env[0] = EnvVar{Key: "K1", Value: "V1"}
	env[1] = EnvVar{Key: "K2", Value: "V2"}
	req := &JobRequest{
		Image: "ubuntu:14.04",
		Cmd:   []string{"env"},
		Env:   env,
	}
	res := req.Execute(testContainerEngine)
	assertNil(t, res.Error)
	assertStringContains(t, "K1=V1\n", string(res.Stdout))
	assertStringContains(t, "K2=V2\n", string(res.Stdout))
	assertString(t, "", string(res.Stderr))
}

func TestStdin(t *testing.T) {
	req := &JobRequest{
		Image: "ubuntu:14.04",
		Cmd:   []string{"cat"},
		Stdin: []byte("hello world\n"),
	}
	res := req.Execute(testContainerEngine)
	assertNil(t, res.Error)
	assertString(t, "hello world\n", string(res.Stdout))
	assertString(t, "", string(res.Stderr))
}

func TestStderr(t *testing.T) {
	req := &JobRequest{
		Image: "ubuntu:14.04",
		Cmd:   []string{"bash", "-c", "echo hello world >&2"},
	}
	res := req.Execute(testContainerEngine)
	assertNil(t, res.Error)
	assertString(t, "", string(res.Stdout))
	assertString(t, "hello world\n", string(res.Stderr))
}

func TestTimeout(t *testing.T) {
	req := &JobRequest{
		Image:   "ubuntu:14.04",
		Cmd:     []string{"sleep", "2000"},
		Timeout: 1 * time.Millisecond,
	}
	res := req.Execute(testContainerEngine)
	assertString(t, errTimeout.Error(), res.Error.Error())
}
