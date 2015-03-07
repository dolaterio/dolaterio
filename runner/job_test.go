package runner

import "testing"

func TestSimpleJob(t *testing.T) {
	req := &JobRequest{
		Image: "busybox",
		Cmd:   []string{"echo", "\"hello world\""},
	}
	res := req.Execute()
	stdout := string(res.Stdout)
	if stdout != "hello world" {
		t.Errorf("Expected \"hello world\" in the stdout, got \"%v\"", stdout)
	}
}
