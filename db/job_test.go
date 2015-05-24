package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreJob(t *testing.T) {
	c, err := NewConnection()
	assert.Nil(t, err)

	job1 := &Job{
		Stdin: "Hello",
	}
	err = job1.Store(c)
	assert.Nil(t, err)
	assert.NotNil(t, job1.ID)

	job2, err := GetJob(c, job1.ID)
	assert.Nil(t, err)
	assert.Equal(t, job1.ID, job2.ID)
	assert.Equal(t, job1.Stdin, job2.Stdin)
}
