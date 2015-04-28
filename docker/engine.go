package docker

import (
	"time"

	"github.com/dolaterio/dolaterio/db"
	"github.com/fsouza/go-dockerclient"
)

// Engine is the engine to process jobs on docker
type Engine struct {
	Timeout time.Duration
	client  *docker.Client
}

// NewEngine initiates and returns a new docker engine.
func NewEngine(config *EngineConfig) (*Engine, error) {
	config.defaults()

	var c *docker.Client
	var err error

	if config.CertPath == "" {
		c, err = docker.NewClient(config.Host)
	} else {
		c, err = docker.NewTLSClient(
			config.Host,
			config.CertPath+"/cert.pem",
			config.CertPath+"/key.pem",
			config.CertPath+"/ca.pem")
	}
	if err != nil {
		return nil, err
	}

	return &Engine{
		client:  c,
		Timeout: config.Timeout,
	}, nil
}

// BuildContainer builds a Container to process the current request
func (engine *Engine) BuildContainer(job *db.Job) (*Container, error) {
	var err error
	// err = engine.client.PullImage(docker.PullImageOptions{
	//  Repository: job.Image,
	// }, docker.AuthConfiguration{})
	// if err != nil {
	//  return nil, err
	// }

	envVars := make([]string, len(job.Env))
	idx := 0
	for k, v := range job.Env {
		envVars[idx] = k + "=" + v
		idx++
	}
	c, err := engine.client.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:      job.DockerImage,
			Env:        envVars,
			Memory:     128 * 1024 * 1024, // 128 MB
			MemorySwap: 0,
			StdinOnce:  true,
			OpenStdin:  true,
			Cmd:        job.Cmd,
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
