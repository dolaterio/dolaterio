package runner

import (
	"bytes"
	"errors"
	"os"
	"time"

	"github.com/fsouza/go-dockerclient"
)

// DockerContainerEngine is the engine to process jobs on docker
type DockerContainerEngine struct{}

// DockerContainer is a data struct representing the container status
type DockerContainer struct {
	containerID string
	stdin       []byte
	stdout      []byte
}

var (
	errTimeout   = errors.New("timeout")
	dockerClient *docker.Client
)

func (*DockerContainerEngine) connect() error {
	if dockerClient != nil {
		return nil
	}
	var err error
	dockerClient, err = docker.NewTLSClient(
		os.Getenv("DOCKER_HOST"),
		os.Getenv("DOCKER_CERT_PATH")+"/cert.pem",
		os.Getenv("DOCKER_CERT_PATH")+"/key.pem",
		os.Getenv("DOCKER_CERT_PATH")+"/ca.pem")
	return err
}

// Run uses the docker engine to run a job
func (engine *DockerContainerEngine) Run(req *JobRequest) (container, error) {
	var err error
	err = engine.connect()
	if err != nil {
		return nil, err
	}

	container, err := engine.buildContainer(req)
	if err != nil {
		return nil, err
	}
	defer container.remove()

	err = container.attachStdin()
	if err != nil {
		return nil, err
	}

	err = container.wait()
	if err != nil {
		return nil, err
	}

	err = container.fetchOutput()
	if err != nil {
		return nil, err
	}

	return container, nil
}

// Stdout returns the stdout of the container
func (container *DockerContainer) Stdout() []byte {
	return container.stdout
}

func (engine *DockerContainerEngine) buildContainer(req *JobRequest) (*DockerContainer, error) {
	var err error
	c, err := dockerClient.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:      req.Image,
			Env:        req.Env.StringArray(),
			Memory:     128 * 1024 * 1024, // 128 MB
			MemorySwap: 0,
			StdinOnce:  true,
			OpenStdin:  true,
			Cmd:        req.Cmd,
		},
	})
	if err != nil {
		return nil, err
	}

	err = dockerClient.StartContainer(c.ID, nil)
	if err != nil {
		return nil, err
	}

	res := &DockerContainer{
		containerID: c.ID,
		stdin:       req.Stdin,
	}

	return res, nil
}

func (container *DockerContainer) attachStdin() error {
	return dockerClient.AttachToContainer(docker.AttachToContainerOptions{
		Container:   container.containerID,
		InputStream: bytes.NewBuffer(container.stdin),
		Stdin:       true,
		Stream:      true,
	})
}

func (container *DockerContainer) remove() error {
	return dockerClient.RemoveContainer(docker.RemoveContainerOptions{
		ID:    container.containerID,
		Force: true,
	})
}

func (container *DockerContainer) fetchOutput() error {
	stdout := new(bytes.Buffer)

	err := dockerClient.Logs(docker.LogsOptions{
		Container:    container.containerID,
		Stdout:       true,
		OutputStream: stdout,
		Tail:         "all",
	})
	if err != nil {
		return err
	}
	container.stdout = stdout.Bytes()
	return nil
}

func (container *DockerContainer) wait() error {
	done := make(chan int)
	errChn := make(chan error)

	go func() {
		_, err := dockerClient.WaitContainer(container.containerID)
		if err != nil {
			errChn <- err
		} else {
			done <- 1
		}
	}()

	select {
	case <-done:
		return nil
	case err := <-errChn:
		return err
	case <-time.After(30 * time.Second):
		return errTimeout
	}
}
