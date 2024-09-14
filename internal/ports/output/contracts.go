package output

import (
	"context"
	"github.com/isaquesb/meli-url-shortener/internal/events"
)

type Dispatcher interface {
	Dispatch(ctx context.Context, msg events.Event) error
	Close()
}
