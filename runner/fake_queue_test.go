package runner

import "github.com/dolaterio/dolaterio/queue"

type fakeQueue struct {
	ids chan string
}

func newFakeQueue() queue.Queue {
	return &fakeQueue{
		ids: make(chan string, 50),
	}
}

func (q *fakeQueue) Enqueue(m *queue.Message) error {
	q.ids <- m.JobID
	return nil
}

func (q *fakeQueue) Dequeue() (*queue.Message, error) {
	m := &queue.Message{
		JobID: <-q.ids,
	}
	return m, nil
}

func (q *fakeQueue) Close() error {
	return nil
}

func (q *fakeQueue) Empty() error {
	return nil
}
