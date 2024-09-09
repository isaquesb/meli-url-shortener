package events

import "context"

//go:generate mockgen -source=contracts.go

type Router interface {
	From(topic string, handler *Handler)
}

type Consumer interface {
	Start(context.Context) error
	GetRouter() Router
}
