package queue

import (
	"fmt"

	"gopkg.in/redis.v2"

	"github.com/Sirupsen/logrus"
	"github.com/dolaterio/dolaterio/core"
)

type redisQueue struct {
	client   *redis.Client
	queueKey string
}

var (
	log = logrus.WithFields(logrus.Fields{
		"package": "queue",
	})
)

// NewRedisQueue returns a redis-backed queue
func NewRedisQueue() (Queue, error) {
	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr: fmt.Sprintf(
			"%v:%v", core.Config.RedisIP, core.Config.RedisPort),
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
	log.WithField("jobID", message.JobID).Info("Queuing message")
	cmd := q.client.RPush(q.queueKey, message.JobID)
	if cmd.Err() != nil {
		log.WithField("jobID", message.JobID).
			WithField("err", cmd.Err()).
			Error("Error queuing message")
		return cmd.Err()
	}
	return nil
}

func (q *redisQueue) Dequeue() (*Message, error) {
	log.Info("Dequeuing message")
	cmd := q.client.BLPop(0, q.queueKey)
	if cmd.Err() != nil {
		log.WithField("err", cmd.Err()).
			Error("Error dequeuing message")
		return nil, cmd.Err()
	}

	return &Message{
		JobID: cmd.Val()[1],
	}, nil
}

func (q *redisQueue) Close() error {
	return q.client.Close()
}
