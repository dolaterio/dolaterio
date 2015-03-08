package runner

import (
	"strings"
	"testing"
)

func assertString(t *testing.T, s1, s2 string) {
	if s1 != s2 {
		t.Errorf("Expected \"%s\", got \"%v\"", s1, s2)
	}
}
func assertStringContains(t *testing.T, s1, s2 string) {
	if strings.Contains(s1, s2) {
		t.Errorf("Expected \"%s\" to contain \"%v\"", s2, s1)
	}
}

func assertNil(t *testing.T, v interface{}) {
	if v != nil {
		t.Errorf("Expected \"%v\" to be nil", v)
	}
}

func TestEcho(t *testing.T) {
	req := &JobRequest{
		Image: "busybox",
		Cmd:   []string{"echo", "hello world"},
	}
	res, err := req.Execute(testContainerEngine)
	assertNil(t, err)
	assertString(t, "hello world\n", string(res.Stdout))
}

func TestEnv(t *testing.T) {
	env := make(EnvVars, 2)
	env[0] = EnvVar{Key: "K1", Value: "V1"}
	env[1] = EnvVar{Key: "K2", Value: "V2"}
	req := &JobRequest{
		Image: "busybox",
		Cmd:   []string{"env"},
		Env:   env,
	}
	res, err := req.Execute(testContainerEngine)
	assertNil(t, err)
	assertStringContains(t, "K1=V1\n", string(res.Stdout))
	assertStringContains(t, "K2=V2\n", string(res.Stdout))
}

func TestStdin(t *testing.T) {
	req := &JobRequest{
		Image: "busybox",
		Cmd:   []string{"cat"},
		Stdin: []byte("hello world\n"),
	}
	res, err := req.Execute(testContainerEngine)
	assertNil(t, err)
	assertString(t, "hello world\n", string(res.Stdout))
}
