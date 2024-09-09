package worker

import (
	"context"
	"github.com/isaquesb/meli-url-shortener/config"
	inputevents "github.com/isaquesb/meli-url-shortener/internal/ports/input/events"
)

func Consume(_ context.Context, consumer inputevents.Consumer) {
	container := config.GetApp()
	consumer.GetRouter().From(container.Events["urls.created"], &inputevents.Handler{
		Name:   "PersistUrl",
		Handle: PersistUrl,
	})

	err := consumer.Start(context.Background())
	if err != nil {
		panic(err)
	}
}
