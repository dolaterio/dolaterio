package queue

// Queue is the interface for a message queue
type Queue interface {
	Enqueue(*Message) error
	Dequeue() (*Message, error)
	Close() error
}

// Message are the messages sent trough the message queue
type Message struct {
	JobID string
}
