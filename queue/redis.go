package queue

import (
	"os"

	"gopkg.in/redis.v2"
)

type redisQueue struct {
	client   *redis.Client
	queueKey string
}

// NewRedisQueue returns a redis-backed queue
func NewRedisQueue() (Queue, error) {
	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    os.Getenv("REDIS_ADDRESS"),
	})

	cmd := client.Ping()
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	return &redisQueue{
		client:   client,
		queueKey: "jobs",
	}, nil
}

func (q *redisQueue) Empty() error {
	cmd := q.client.Del(q.queueKey)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (q *redisQueue) Enqueue(message *Message) error {
	cmd := q.client.RPush(q.queueKey, message.JobID)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (q *redisQueue) Dequeue() (*Message, error) {
	cmd := q.client.BLPop(0, q.queueKey)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	return &Message{
		JobID: cmd.Val()[1],
	}, nil
}

func (q *redisQueue) Close() error {
	return q.client.Close()
}
