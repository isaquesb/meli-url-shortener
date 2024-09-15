package urls

import (
	"fmt"
	"github.com/isaquesb/meli-url-shortener/internal/app"
	inputevents "github.com/isaquesb/meli-url-shortener/internal/events"
	"strings"
)

func PersistHandler(m *inputevents.Message) error {
	container := app.GetApp()

	dispatcher := container.Worker.Dispatcher.Get()
	return dispatcher.Dispatch(container.Ctx, m.Event)
}

func EventSubscriberFor(name string) *inputevents.Subscriber {
	handlerName := fmt.Sprintf("%sUrlHandler", strings.Replace(name, "urls.", "", -1))
	return &inputevents.Subscriber{
		Name:       strings.ToUpper(handlerName[0:1]) + handlerName[1:],
		ParseEvent: EventParserFor(name),
		Handler:    PersistHandler,
	}
}
