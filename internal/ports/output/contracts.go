package output

import (
	"context"
	"github.com/isaquesb/url-shortener/internal/events"
)

type Dispatcher interface {
	Dispatch(ctx context.Context, msg events.Event) error
	Close()
}

type Listen interface {
	Listen(context.Context)
}

type UrlRepository interface {
	UrlFromShort(context.Context, string) (string, error)
	StatsFromShort(context.Context, string) (map[string]interface{}, error)
}
