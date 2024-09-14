package worker

import (
	"context"
	"github.com/isaquesb/meli-url-shortener/internal/app"
	"github.com/isaquesb/meli-url-shortener/internal/events"
	inputevents "github.com/isaquesb/meli-url-shortener/internal/ports/input/events"
	"github.com/isaquesb/meli-url-shortener/internal/urls"
)

func Consume(_ context.Context, consumer inputevents.Consumer) {
	container := app.GetApp()
	consumer.GetRouter().From(container.Events[urls.Created], &events.Handler{
		Id:         "UrlCreateHandler",
		ParseEvent: func() (events.Event, error) { return urls.MakeEvent(urls.Created) },
		Handle:     urls.CreateHandler,
	})

	err := consumer.Start(context.Background())
	if err != nil {
		panic(err)
	}
}
