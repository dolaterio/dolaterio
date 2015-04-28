package runner

// import (
// 	"testing"
// 	"time"

// 	"github.com/dolaterio/dolaterio/db"
// 	"github.com/dolaterio/dolaterio/docker"
// 	"github.com/stretchr/testify/assert"
// )

// func TestEcho(t *testing.T) {
// 	job := &db.Job{
// 		DockerImage: "ubuntu:14.04",
// 		Cmd:         []string{"echo", "hello world"},
// 	}
// 	engine, err := docker.NewEngine(&docker.EngineConfig{})
// 	assert.Nil(t, err)

// 	err = Run(job, engine)
// 	assert.Nil(t, err)
// 	assert.Equal(t, "hello world\n", string(job.Stdout))
// 	assert.Equal(t, "", string(job.Stderr))
// }

// func TestEnv(t *testing.T) {
// 	job := &db.Job{
// 		DockerImage: "ubuntu:14.04",
// 		Cmd:         []string{"env"},
// 		Env:         map[string]string{"K1": "V1", "K2": "V2"},
// 	}
// 	engine, err := docker.NewEngine(&docker.EngineConfig{})
// 	assert.Nil(t, err)

// 	err = Run(job, engine)
// 	assert.Nil(t, err)
// 	assert.Contains(t, job.Stdout, "K1=V1")
// 	assert.Contains(t, job.Stdout, "K2=V2")
// }

// func TestStdin(t *testing.T) {
// 	job := &db.Job{
// 		DockerImage: "ubuntu:14.04",
// 		Cmd:         []string{"cat"},
// 		Stdin:       "hello world\n",
// 	}
// 	engine, err := docker.NewEngine(&docker.EngineConfig{})
// 	assert.Nil(t, err)

// 	err = Run(job, engine)
// 	assert.Nil(t, err)
// 	assert.Equal(t, "hello world\n", job.Stdout)
// }

// func TestStderr(t *testing.T) {
// 	job := &db.Job{
// 		DockerImage: "ubuntu:14.04",
// 		Cmd:         []string{"bash", "-c", "echo hello world >&2"},
// 	}
// 	engine, err := docker.NewEngine(&docker.EngineConfig{})
// 	assert.Nil(t, err)

// 	err = Run(job, engine)
// 	assert.Nil(t, err)
// 	assert.Equal(t, "hello world\n", job.Stderr)
// }

// func TestTimeout(t *testing.T) {
// 	job := &db.Job{
// 		DockerImage: "ubuntu:14.04",
// 		Cmd:         []string{"sleep", "2000"},
// 		Timeout:     1 * time.Millisecond,
// 	}
// 	engine, err := docker.NewEngine(&docker.EngineConfig{})
// 	assert.Nil(t, err)

// 	err = Run(job, engine)
// 	assert.Equal(t, errTimeout.Error(), err.Error())
// }
