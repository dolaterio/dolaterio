package dolaterio

import "time"

// Job models a job
type Job struct {
	ID      string
	Image   string
	Cmd     []string
	Stdin   []byte
	Env     EnvVars
	Timeout time.Duration
	Stdout  []byte
	Stderr  []byte
	Error   error
}
