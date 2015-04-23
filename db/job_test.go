package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreJob(t *testing.T) {
	err := Connect()
	assert.Nil(t, err)

	job1 := &Job{
		DockerImage: "ubuntu:14.04",
	}
	err = job1.Store()
	assert.Nil(t, err)
	assert.NotNil(t, job1.ID)

	job2, err := GetJob(job1.ID)
	assert.Nil(t, err)
	assert.Equal(t, job1.ID, job2.ID)
	assert.Equal(t, job1.DockerImage, job2.DockerImage)
}
