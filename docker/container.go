package docker

import (
	"bytes"

	"github.com/fsouza/go-dockerclient"
)

// Container is a data struct representing the container status
type Container struct {
	engine      *Engine
	containerID string

	StdIn  []byte
	StdOut []byte
	StdErr []byte
}

// AttachStdin sends the StdIn to the container
func (container *Container) AttachStdin() error {
	if container.StdIn == nil {
		return nil
	}
	return container.engine.client.AttachToContainer(docker.AttachToContainerOptions{
		Container:   container.containerID,
		InputStream: bytes.NewBuffer(container.StdIn),
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
	container.StdOut = stdout.Bytes()
	container.StdErr = stderr.Bytes()
	return nil
}

// Remove removes the container from the docker host
func (container *Container) Remove() error {
	return container.engine.client.RemoveContainer(docker.RemoveContainerOptions{
		ID:    container.containerID,
		Force: true,
	})
}
