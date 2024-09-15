package worker

import (
	"context"
	"github.com/isaquesb/meli-url-shortener/internal/app"
	inputevents "github.com/isaquesb/meli-url-shortener/internal/ports/input/events"
	"github.com/isaquesb/meli-url-shortener/internal/urls"
)

func Consume(_ context.Context, consumer inputevents.Consumer) {
	container := app.GetApp()
	router := consumer.GetRouter()
	router.From(container.Topics[urls.Created], urls.EventSubscriberFor(urls.Created))
	router.From(container.Topics[urls.Visited], urls.EventSubscriberFor(urls.Visited))
	router.From(container.Topics[urls.Deleted], urls.EventSubscriberFor(urls.Deleted))

	err := consumer.Start(context.Background())
	if err != nil {
		panic(err)
	}
}
