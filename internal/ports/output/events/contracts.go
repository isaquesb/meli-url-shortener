package events

import "context"

const TypePlain = "plain"
const TypeJson = "json"

type Message struct {
	ContentType string
	Key         []byte
	Body        []byte
}

//go:generate mockgen -source=contracts.go

type Dispatcher interface {
	Dispatch(ctx context.Context, to string, msg Message) error
	Close()
}
