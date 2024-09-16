package watermill

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	inputEvents "github.com/isaquesb/url-shortener/internal/ports/input/events"
)

func NewConsumer(logger watermill.LoggerAdapter, subscriber message.Subscriber) *Consumer {
	return &Consumer{
		Logger: logger,
		Router: &Router{
			MessageRouter: NewMessageRouter(logger),
		},
		Subscriber: subscriber,
	}
}

type Consumer struct {
	Logger     watermill.LoggerAdapter
	Router     *Router
	Subscriber message.Subscriber
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
