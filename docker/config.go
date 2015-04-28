package docker

import (
	"os"
	"time"
)

// EngineConfig is the configuration definition to initiate a new engine
type EngineConfig struct {
	Host     string
	CertPath string
	Timeout  time.Duration
}

func (e *EngineConfig) defaults() {
	if e.Host == "" {
		e.Host = os.Getenv("DOCKER_HOST")
	}
	if e.Host == "" {
		e.Host = "unix:///var/run/docker.sock"
	}

	if e.CertPath == "" {
		e.CertPath = os.Getenv("DOCKER_CERT_PATH")
	}

	if e.Timeout == 0 {
		e.Timeout = 30 * time.Second
	}
}
