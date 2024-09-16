package watermill

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/isaquesb/url-shortener/internal/events"
)

func NewMessageRouter(logger watermill.LoggerAdapter) *message.Router {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}
	router.AddPlugin(plugin.SignalsHandler)

	return router
}

type Route struct {
	Topic      string
	Subscriber *events.Subscriber
}

type Router struct {
	MessageRouter *message.Router
	Routes        []*Route
}

func (r *Router) From(topic string, subscriber *events.Subscriber) {
	r.Routes = append(r.Routes, &Route{Topic: topic, Subscriber: subscriber})
}

func (r *Router) RoutedSubscriber(subscriber *events.Subscriber) func(m *message.Message) error {
	return func(m *message.Message) error {
		evt, err := subscriber.ParseEvent()
		if err != nil {
			return err
		}
		decoded := &events.Envelop{Event: evt}
		if err := json.Unmarshal(m.Payload, decoded); err != nil {
			return err
		}
		return subscriber.Handler(&events.Message{
			Uuid:  m.UUID,
			Name:  decoded.Name,
			Event: decoded.Event,
		})
	}
}
