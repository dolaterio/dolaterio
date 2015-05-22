package queue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEnqueueDequeue(t *testing.T) {
	queue, err := NewRedisQueue()
	assert.Nil(t, err)
	defer queue.Close()

	err = queue.Empty()
	assert.Nil(t, err)

	jobID := "MyJob"
	err = queue.Enqueue(&Message{JobID: jobID})
	assert.Nil(t, err)

	message, err := queue.Dequeue()
	assert.Nil(t, err)
	assert.Equal(t, jobID, message.JobID)
}

func TestDequeueWaits(t *testing.T) {
	queue, err := NewRedisQueue()
	assert.Nil(t, err)
	defer queue.Close()

	err = queue.Empty()
	assert.Nil(t, err)

	jobID := "MyJob"

	go func() {
		begin := time.Now()
		message, err := queue.Dequeue()
		assert.Nil(t, err)
		assert.Equal(t, jobID, message.JobID)
		took := time.Since(begin)
		assert.True(t, took > 50*time.Millisecond, "Expected to wait for at least 50ms")
	}()

	time.Sleep(50 * time.Millisecond)
	err = queue.Enqueue(&Message{JobID: jobID})
	assert.Nil(t, err)
	time.Sleep(50 * time.Millisecond)
}

func TestOnlyOneConsumerGetsAMessage(t *testing.T) {
	queue, err := NewRedisQueue()
	assert.Nil(t, err)
	defer queue.Close()

	err = queue.Empty()
	assert.Nil(t, err)

	jobID := "MyJob"

	count := 0

	for i := 0; i < 10; i++ {
		go func() {
			queue.Dequeue()
			count++
		}()
	}

	err = queue.Enqueue(&Message{JobID: jobID})
	assert.Nil(t, err)
	time.Sleep(5 * time.Millisecond)

	assert.Equal(t, 1, count)
}
