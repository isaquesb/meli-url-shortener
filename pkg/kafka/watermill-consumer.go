package kafka

import (
	"context"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill"
	watermillKafka "github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/isaquesb/meli-url-shortener/internal/events"
	inputEvents "github.com/isaquesb/meli-url-shortener/internal/ports/input/events"
)

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

type Consumer struct {
	Logger     watermill.LoggerAdapter
	Router     *Router
	Subscriber message.Subscriber
}

func NewMessageRouter(logger watermill.LoggerAdapter) *message.Router {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}
	router.AddPlugin(plugin.SignalsHandler)

	return router
}

func NewLogger(debug, trace bool) watermill.LoggerAdapter {
	return watermill.NewStdLogger(debug, trace)
}

func NewSubscriber(logger watermill.LoggerAdapter, consumerGroup string, brokers []string) (message.Subscriber, error) {
	sub, err := watermillKafka.NewSubscriber(watermillKafka.SubscriberConfig{
		Brokers:       brokers,
		Unmarshaler:   watermillKafka.DefaultMarshaler{},
		ConsumerGroup: consumerGroup,
	}, logger)

	if err != nil {
		return nil, err
	}
	return sub, nil
}

func NewConsumer(logger watermill.LoggerAdapter, subscriber message.Subscriber) *Consumer {
	return &Consumer{
		Logger: logger,
		Router: &Router{
			MessageRouter: NewMessageRouter(logger),
		},
		Subscriber: subscriber,
	}
}

func (c *Consumer) GetRouter() inputEvents.Router {
	return c.Router
}

func (c *Consumer) Start(ctx context.Context) error {
	for _, route := range c.Router.Routes {
		c.Router.MessageRouter.AddNoPublisherHandler(
			route.Subscriber.Name,
			route.Topic,
			c.Subscriber,
			c.Router.RoutedSubscriber(route.Subscriber),
		)
	}

	return c.Router.MessageRouter.Run(ctx)
}
