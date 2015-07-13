package docker

import (
	"time"

	// Engine is the engine to process jobs on docker
	"github.com/dolaterio/dolaterio/core"
	"github.com/dolaterio/dolaterio/db"
	"github.com/fsouza/go-dockerclient"
)

type Engine struct {
	Timeout time.Duration
	client  *docker.Client
}

// NewEngine initiates and returns a new docker engine.
func NewEngine() (*Engine, error) {
	var c *docker.Client
	var err error

	if core.Config.DockerCertPath == "" {
		c, err = docker.NewClient(core.Config.DockerHost)
	} else {
		c, err = docker.NewTLSClient(
			core.Config.DockerHost,
			core.Config.DockerCertPath+"/cert.pem",
			core.Config.DockerCertPath+"/key.pem",
			core.Config.DockerCertPath+"/ca.pem")
	}
	if err != nil {
		return nil, err
	}

	return &Engine{
		client:  c,
		Timeout: core.Config.TaskTimeout,
	}, nil
}

// BuildContainer builds a Container to process the current request
func (engine *Engine) BuildContainer(job *db.Job) (*Container, error) {
	var err error
	err = engine.client.PullImage(docker.PullImageOptions{
		Repository: job.Worker.DockerImage,
	}, docker.AuthConfiguration{})
	if err != nil {
		return nil, err
	}

	mergedVars := map[string]string{}
	if job.Worker.Env != nil {
		for k, v := range job.Worker.Env {
			mergedVars[k] = v
		}
	}
	if job.Env != nil {
		for k, v := range job.Env {
			mergedVars[k] = v
		}
	}

	envVars := make([]string, len(mergedVars))
	idx := 0
	for k, v := range mergedVars {
		envVars[idx] = k + "=" + v
		idx++
	}
	c, err := engine.client.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:      job.Worker.DockerImage,
			Env:        envVars,
			Memory:     128 * 1024 * 1024, // 128 MB
			MemorySwap: 0,
			StdinOnce:  true,
			OpenStdin:  true,
			Cmd:        job.Worker.Cmd,
		},
	})
	if err != nil {
		return nil, err
	}

	err = engine.client.StartContainer(c.ID, nil)
	if err != nil {
		return nil, err
	}

	res := &Container{
		engine:      engine,
		containerID: c.ID,
		StdIn:       []byte(job.Stdin),
	}

	return res, nil
}
