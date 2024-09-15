package events

import (
	"context"
	internalEvents "github.com/isaquesb/url-shortener/internal/events"
)

type Router interface {
	From(string, *internalEvents.Subscriber)
}

type Consumer interface {
	Start(context.Context) error
	GetRouter() Router
}
