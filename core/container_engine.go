package dolaterio

import (
	"bytes"
	"errors"
	"os"
	"time"

	"github.com/fsouza/go-dockerclient"
)

// ContainerEngine is the engine to process jobs on docker
type ContainerEngine struct {
	client *docker.Client
}

// Container is a data struct representing the container status
type Container struct {
	engine      *ContainerEngine
	containerID string
	stdin       []byte
	stdout      []byte
	stderr      []byte
}

var (
	errTimeout = errors.New("timeout")
)

// Connect connects to the docker host and sets the client
func (engine *ContainerEngine) Connect() error {
	var c *docker.Client
	var err error

	host := os.Getenv("DOCKER_HOST")
	if host == "" {
		host = "unix:///var/run/docker.sock"
	}

	if os.Getenv("DOCKER_CERT_PATH") == "" {
		c, err = docker.NewClient(host)
	} else {
		c, err = docker.NewTLSClient(
			host,
			os.Getenv("DOCKER_CERT_PATH")+"/cert.pem",
			os.Getenv("DOCKER_CERT_PATH")+"/key.pem",
			os.Getenv("DOCKER_CERT_PATH")+"/ca.pem")
	}
	engine.client = c
	return err
}

// Timeout returns the default timeout
func (engine *ContainerEngine) Timeout() time.Duration {
	return 30 * time.Second
}

// BuildContainer builds a Container to process the current request
func (engine *ContainerEngine) BuildContainer(req *JobRequest) (*Container, error) {
	var err error
	// err = engine.client.PullImage(docker.PullImageOptions{
	// 	Repository: req.Image,
	// }, docker.AuthConfiguration{})
	// if err != nil {
	// 	return nil, err
	// }
	c, err := engine.client.CreateContainer(docker.CreateContainerOptions{
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

	err = engine.client.StartContainer(c.ID, nil)
	if err != nil {
		return nil, err
	}

	res := &Container{
		engine:      engine,
		containerID: c.ID,
		stdin:       req.Stdin,
	}

	return res, nil
}

// AttachStdin sends the stdin to the container
func (container *Container) AttachStdin() error {
	if container.stdin == nil {
		return nil
	}
	return container.engine.client.AttachToContainer(docker.AttachToContainerOptions{
		Container:   container.containerID,
		InputStream: bytes.NewBuffer(container.stdin),
		Stdin:       true,
		Stream:      true,
	})
}

// Wait waits for the docker container to be done (or timeout in 30s)
func (container *Container) Wait() error {
	_, err := container.engine.client.WaitContainer(container.containerID)
	return err
}

// FetchOutput retrieves the outputs from the container
func (container *Container) FetchOutput() error {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	err := container.engine.client.Logs(docker.LogsOptions{
		Container:    container.containerID,
		Stdout:       true,
		Stderr:       true,
		OutputStream: stdout,
		ErrorStream:  stderr,
		Tail:         "all",
	})
	if err != nil {
		return err
	}
	container.stdout = stdout.Bytes()
	container.stderr = stderr.Bytes()
	return nil
}

// Remove removes the container from the docker host
func (container *Container) Remove() error {
	return container.engine.client.RemoveContainer(docker.RemoveContainerOptions{
		ID:    container.containerID,
		Force: true,
	})
}

// Stdout returns the stdout of the container
func (container *Container) Stdout() []byte {
	return container.stdout
}

// Stderr returns the stderr of the container
func (container *Container) Stderr() []byte {
	return container.stderr
}
