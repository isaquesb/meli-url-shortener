package worker

import (
	"fmt"
	"github.com/isaquesb/meli-url-shortener/config"
	inputevents "github.com/isaquesb/meli-url-shortener/internal/ports/input/events"
	"github.com/isaquesb/meli-url-shortener/internal/ports/output/events"
	"github.com/isaquesb/meli-url-shortener/pkg/logger"
)

func PersistUrl(m *inputevents.Message) error {
	strPayload := string(m.Payload)
	short := strPayload[0:6]
	url := strPayload[6:]

	container := config.GetApp()

	dispatcher := container.Worker.GetDispatcher()
	err := dispatcher.Dispatch(container.Ctx, "urls", events.Message{
		ContentType: events.TypePlain,
		Key:         []byte(short),
		Body:        []byte(url),
	})

	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("persisted short: %s, url: %s", short, url))

	return nil
}
