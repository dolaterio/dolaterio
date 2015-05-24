package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreWorker(t *testing.T) {
	c, err := NewConnection()
	assert.Nil(t, err)

	worker1 := &Worker{
		DockerImage: "ubuntu:14.04",
	}
	err = worker1.Store(c)
	assert.Nil(t, err)
	assert.NotNil(t, worker1.ID)

	worker2, err := GetWorker(c, worker1.ID)
	assert.Nil(t, err)
	assert.Equal(t, worker1.ID, worker2.ID)
	assert.Equal(t, worker1.DockerImage, worker2.DockerImage)
}
