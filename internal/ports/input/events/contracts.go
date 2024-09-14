package events

import (
	"context"
	events2 "github.com/isaquesb/meli-url-shortener/internal/events"
)

type Router interface {
	From(topic string, handler *events2.Handler)
}

type Consumer interface {
	Start(context.Context) error
	GetRouter() Router
}
