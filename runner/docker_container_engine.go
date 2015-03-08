package runner

import (
	"bytes"
	"os"

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
func (engine *DockerContainerEngine) Run(image string, cmd []string, env EnvVars) (container, error) {
	engine.connect()

	createContainerOpts := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:      image,
			Env:        env.StringArray(),
			Memory:     128 * 1024 * 1024, // 128 MB
			MemorySwap: 0,
			StdinOnce:  true,
			OpenStdin:  true,
		},
	}
	c, err := dockerClient.CreateContainer(createContainerOpts)
	if err != nil {
		return nil, err
	}
	defer dockerClient.RemoveContainer(docker.RemoveContainerOptions{ID: c.ID, Force: true})
	res := &DockerContainer{
		containerID: c.ID,
	}

	if err := dockerClient.StartContainer(c.ID, nil); err != nil {
		return nil, err
	}
	//
	// attachContainerOpts := docker.AttachToContainerOptions{
	//   Container:   c.ID,
	//   InputStream: bytes.NewBuffer(job.Payload),
	//   Stdin:       true,
	//   Stream:      true,
	// }
	// if err := client.AttachToContainer(attachContainerOpts); err != nil {
	//   return &JobResponse{Error: err}
	// }

	_, err = dockerClient.WaitContainer(c.ID)
	if err != nil {
		return nil, err
	}

	c, err = dockerClient.InspectContainer(c.ID)
	if err != nil {
		return nil, err
	}

	stdout := new(bytes.Buffer)

	logs := docker.LogsOptions{
		Container:    c.ID,
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
