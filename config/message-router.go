package config

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	pkgWatermill "github.com/isaquesb/meli-url-shortener/pkg/watermill"
)

func MessageRouterLogger() watermill.LoggerAdapter {
	return pkgWatermill.NewLogger(
		GetEnv("APP_DEBUG", "false") == "true",
		GetEnv("APP_TRACE", "false") == "true",
	)
}

func MessageRouter(logger watermill.LoggerAdapter) *message.Router {
	return pkgWatermill.NewRouter(logger)
}

func KafkaSubscriber(logger watermill.LoggerAdapter) (message.Subscriber, error) {
	sub, err := kafka.NewSubscriber(kafka.SubscriberConfig{
		Brokers:       []string{GetEnv("KAFKA_BROKER", "localhost:9092")},
		Unmarshaler:   kafka.DefaultMarshaler{},
		ConsumerGroup: GetEnv("KAFKA_CONSUMER_GROUP", "url-shortener"),
	}, logger)

	if err != nil {
		return nil, err
	}
	return sub, nil
}
