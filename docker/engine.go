package docker

import (
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dolaterio/dolaterio/core"
	"github.com/dolaterio/dolaterio/db"
	"github.com/fsouza/go-dockerclient"
)

// Engine is the engine to process jobs on docker
type Engine struct {
	Timeout  time.Duration
	SkipPull bool
	client   *docker.Client
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

// ValidImage builds a Container to process the current request
func (engine *Engine) ValidImage(imageName string) (bool, error) {
	log.WithField("image", imageName).Info("Validating docker image")
	imageNameWithoutTag := strings.SplitN(imageName, ":", 2)[0]
	images, err := engine.client.SearchImages(imageNameWithoutTag)
	if err != nil {
		log.WithFields(logrus.Fields{"image": imageNameWithoutTag, "err": err}).
			Error("Error validating image")
		return false, err
	}
	log.WithFields(logrus.Fields{"image": imageNameWithoutTag, "found": len(images)}).
		Debug("Images found")
	return len(images) > 0, nil
}

// BuildContainer builds a Container to process the current request
func (engine *Engine) BuildContainer(job *db.Job) (*Container, error) {
	logFields := logrus.Fields{"jobID": job.ID}
	var err error

	log.WithFields(logFields).Info("Building container")

	if !engine.SkipPull {
		log.WithFields(logFields).Info("Pulling docker image")
		err = engine.client.PullImage(docker.PullImageOptions{
			Repository: job.Worker.DockerImage,
		}, docker.AuthConfiguration{})
		if err != nil {
			log.WithFields(logFields).WithField("err", err).
				Error("Error fetching image")
			return nil, err
		}
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
		log.WithFields(logFields).WithField("err", err).
			Error("Error creating the container")
		return nil, err
	}

	err = engine.client.StartContainer(c.ID, nil)
	if err != nil {
		log.WithFields(logFields).WithField("err", err).
			Error("Error starting the container")
		return nil, err
	}

	res := &Container{
		engine:      engine,
		containerID: c.ID,
		StdIn:       []byte(job.Stdin),
	}

	return res, nil
}
