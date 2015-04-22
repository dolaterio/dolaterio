package dolaterio

type Queue interface {
	Pop() (error, *QueueMessage)
}

type QueueMessage struct {
	JobID int
}
