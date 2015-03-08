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

	createContainerOpts := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:      req.Image,
			Env:        req.Env.StringArray(),
			Memory:     128 * 1024 * 1024, // 128 MB
			MemorySwap: 0,
			StdinOnce:  true,
			OpenStdin:  true,
			Cmd:        req.Cmd,
		},
	}

	container, err := dockerClient.CreateContainer(createContainerOpts)
	if err != nil {
		return nil, err
	}
	defer dockerClient.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID, Force: true})
	res := &DockerContainer{
		containerID: container.ID,
	}

	if err := dockerClient.StartContainer(container.ID, nil); err != nil {
		return nil, err
	}

	attachContainerOpts := docker.AttachToContainerOptions{
		Container:   container.ID,
		InputStream: bytes.NewBuffer(req.Stdin),
		Stdin:       true,
		Stream:      true,
	}
	if err := dockerClient.AttachToContainer(attachContainerOpts); err != nil {
		return nil, err
	}

	err = res.wait()
	if err != nil {
		return nil, err
	}

	container, err = dockerClient.InspectContainer(container.ID)
	if err != nil {
		return nil, err
	}

	stdout := new(bytes.Buffer)

	logs := docker.LogsOptions{
		Container:    container.ID,
		Stdout:       true,
		OutputStream: stdout,
		Tail:         "all",
	}
	if err := dockerClient.Logs(logs); err != nil {
		return nil, err
	}
	res.stdout = stdout.Bytes()

	return res, nil
}

// Stdout returns the stdout of the container
func (container *DockerContainer) Stdout() []byte {
	return container.stdout
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
