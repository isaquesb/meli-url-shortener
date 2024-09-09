package events

import "context"

type Router interface {
	From(topic string, handler *Handler)
}

type Consumer interface {
	Start(context.Context) error
	GetRouter() Router
}
