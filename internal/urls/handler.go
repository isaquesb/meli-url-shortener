package urls

import (
	"fmt"
	inputevents "github.com/isaquesb/meli-url-shortener/internal/events"
	"github.com/isaquesb/meli-url-shortener/pkg/logger"
)

func CreateHandler(m *inputevents.Message) error {
	evt := m.Event.(*CreateEvent)

	/*container := config.GetApp()

	dispatcher := container.Worker.GetDispatcher()
	err := dispatcher.Dispatch(container.Ctx, events.Event{
		Operation: events.OpCreate,
		Key:       []byte(short),
		Body:      []byte(url),
	})

	if err != nil {
		return err
	}*/

	logger.Info(fmt.Sprintf("persisted short: %s, url: %s, createdAt: %v", evt.ShortCode, evt.Url, evt.CreatedAt))

	return nil
}
