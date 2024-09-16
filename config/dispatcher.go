package config

import (
	"github.com/isaquesb/url-shortener/internal/ports/output"
	"github.com/isaquesb/url-shortener/pkg/kafka"
	"github.com/isaquesb/url-shortener/pkg/rabbitmq"
)

func GetAppDispatcher() (output.Dispatcher, error) {
	switch GetEnv("APP_DISPATCHER", "kafka") {
	case "rabbitmq":
		return getRabbitMqDispatcher()
	default:
		return getKafkaDispatcher()
	}
}

func getKafkaDispatcher() (output.Dispatcher, error) {
	return kafka.NewDispatcher(map[string]interface{}{
		"bootstrap.servers": GetEnv("KAFKA_BROKERS", "kafka:9092"),
	})
}

func getRabbitMqDispatcher() (output.Dispatcher, error) {
	return rabbitmq.NewDispatcher(map[string]interface{}{
		"url": getAmqpURI(),
	})
}
