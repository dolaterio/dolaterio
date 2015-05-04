package queue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEnqueueDequeue(t *testing.T) {
	queue, err := NewRedisQueue()
	assert.Nil(t, err)

	jobID := "MyJob"
	err = queue.Enqueue(&Message{JobID: jobID})
	assert.Nil(t, err)

	message, err := queue.Dequeue()
	assert.Nil(t, err)
	assert.Equal(t, message.JobID, jobID)
}

func TestDequeueWaits(t *testing.T) {
	queue, err := NewRedisQueue()
	assert.Nil(t, err)

	jobID := "MyJob"

	go func() {
		message, err := queue.Dequeue()
		assert.Nil(t, err)
		assert.Equal(t, message.JobID, jobID)
	}()

	time.Sleep(50 * time.Millisecond)
	err = queue.Enqueue(&Message{JobID: jobID})
	assert.Nil(t, err)
}
