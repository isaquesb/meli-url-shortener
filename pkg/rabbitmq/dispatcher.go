package rabbitmq

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/isaquesb/url-shortener/internal/events"
	"github.com/isaquesb/url-shortener/internal/ports/output"
	"github.com/isaquesb/url-shortener/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type Dispatcher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewDispatcher(cfg map[string]interface{}) (output.Dispatcher, error) {
	config := amqp.Config{
		Vhost:      "/",
		Properties: amqp.NewConnectionProperties(),
	}
	config.Properties.SetClientConnectionName("url-shortener")
	conn, err := amqp.Dial(cfg["url"].(string))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		"url-shortener", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)

	if err != nil {
		return nil, err
	}

	return &Dispatcher{
		conn:    conn,
		channel: ch,
	}, nil
}

func (dp *Dispatcher) Close() {
	if err := dp.conn.Close(); err != nil {
		logger.Error("Failed to close connection: %v", err)
	}
	if err := dp.channel.Close(); err != nil {
		logger.Error("Failed to close channel: %v", err)
	}
}

func (dp *Dispatcher) Dispatch(ctx context.Context, msg events.Event) error {
	encoded, err := json.Marshal(events.Envelop{
		Name:  msg.GetName(),
		Event: msg,
	})
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return dp.channel.PublishWithContext(
		ctx,
		msg.GetName(), // Exchange
		msg.GetName(), // Routing key
		false,         // Mandatory
		false,         // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        encoded,
		},
	)
}
