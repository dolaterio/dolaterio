package dolaterio

import "time"

type containerEngine interface {
	Connect() error
	BuildContainer(*JobRequest) (container, error)
	Timeout() time.Duration
}

type container interface {
	AttachStdin() error
	Wait() error
	Remove() error
	FetchOutput() error

	Stdout() []byte
	Stderr() []byte
}
