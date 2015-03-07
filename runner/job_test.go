package runner

import "testing"

func assertString(t *testing.T, s1, s2 string) {
	if s1 != s2 {
		t.Errorf("Expected \"%s\", got \"%v\"", s1, s2)
	}
}

func TestEcho(t *testing.T) {
	req := &JobRequest{
		Image: "busybox",
		Cmd:   []string{"echo", "hello world"},
	}
	res := req.Execute(testContainerEngine)
	assertString(t, "hello world", string(res.Stdout))
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
	res := req.Execute(testContainerEngine)
	assertString(t, "K1=V1\nK2=V2\n", string(res.Stdout))
}
