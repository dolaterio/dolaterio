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
	assert.Equal(t, jobID, message.JobID)
	queue.Close()
}

func TestDequeueWaits(t *testing.T) {
	queue, err := NewRedisQueue()
	assert.Nil(t, err)

	jobID := "MyJob"

	go func() {
		message, err := queue.Dequeue()
		assert.Nil(t, err)
		assert.Equal(t, jobID, message.JobID)
		queue.Close()
	}()

	time.Sleep(50 * time.Millisecond)
	err = queue.Enqueue(&Message{JobID: jobID})
	assert.Nil(t, err)
}

func TestOnlyOneConsumerGetsAMessage(t *testing.T) {
	jobID := "MyJob"

	count := 0

	for i := 0; i < 10; i++ {
		go func() {
			queue, err := NewRedisQueue()
			assert.Nil(t, err)

			_, err = queue.Dequeue()
			count++
		}()
	}

	queue, err := NewRedisQueue()
	assert.Nil(t, err)

	err = queue.Enqueue(&Message{JobID: jobID})
	assert.Nil(t, err)
	time.Sleep(1 * time.Millisecond)

	assert.Equal(t, 1, count)
	queue.Close()
}
