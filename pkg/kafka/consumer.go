package kafka

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	wkafka "github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/isaquesb/meli-url-shortener/internal/ports/input/events"
)

type Route struct {
	Topic   string
	Handler *events.Handler
}

type Router struct {
	MessageRouter *message.Router
	Routes        []*Route
}

func (r *Router) From(topic string, handler *events.Handler) {
	r.Routes = append(r.Routes, &Route{Topic: topic, Handler: handler})
}

func (r *Router) RoutedHandler(handler *events.Handler) func(m *message.Message) error {
	return func(m *message.Message) error {
		msg := &events.Message{
			Uuid:    m.UUID,
			Payload: m.Payload,
		}
		return handler.Handle(msg)
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
	sub, err := wkafka.NewSubscriber(wkafka.SubscriberConfig{
		Brokers:       brokers,
		Unmarshaler:   wkafka.DefaultMarshaler{},
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

func (c *Consumer) GetRouter() events.Router {
	return c.Router
}

func (c *Consumer) Start(ctx context.Context) error {
	for _, route := range c.Router.Routes {
		c.Router.MessageRouter.AddNoPublisherHandler(
			route.Handler.Name,
			route.Topic,
			c.Subscriber,
			c.Router.RoutedHandler(route.Handler),
		)
	}

	return c.Router.MessageRouter.Run(ctx)
}
