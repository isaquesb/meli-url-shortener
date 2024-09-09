package events

import "context"

const OpCreate = "create"

type Message struct {
	Operation string
	Key       []byte
	Body      []byte
}

type Dispatcher interface {
	Dispatch(ctx context.Context, to string, msg Message) error
	Close()
}
