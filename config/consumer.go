package config

import (
	"github.com/ThreeDotsLabs/watermill"
	inputevents "github.com/isaquesb/url-shortener/internal/ports/input/events"
	"github.com/isaquesb/url-shortener/pkg/kafka"
	"github.com/isaquesb/url-shortener/pkg/rabbitmq"
	watermill2 "github.com/isaquesb/url-shortener/pkg/watermill"
)

func GetAppConsumer(logger watermill.LoggerAdapter, group string) (inputevents.Consumer, error) {
	switch GetEnv("APP_DISPATCHER", "kafka") {
	case "rabbitmq":
		return getRabbitMqConsumer(logger, group)
	default:
		return getKafkaConsumer(logger, group)
	}
}

func getKafkaConsumer(logger watermill.LoggerAdapter, group string) (inputevents.Consumer, error) {
	subscriber, err := kafka.NewSubscriber(
		logger,
		group,
		[]string{GetEnv("KAFKA_BROKERS", "kafka:9092")},
	)
	if err != nil {
		return nil, err
	}
	return watermill2.NewConsumer(logger, subscriber), nil
}

func getRabbitMqConsumer(logger watermill.LoggerAdapter, group string) (inputevents.Consumer, error) {
	subscriber, err := rabbitmq.NewSubscriber(
		logger,
		group,
		getAmqpURI(),
	)
	if err != nil {
		return nil, err
	}
	return watermill2.NewConsumer(logger, subscriber), nil
}
