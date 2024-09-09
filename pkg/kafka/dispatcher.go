package kafka

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/isaquesb/meli-url-shortener/internal/ports/output/events"
	"log"
)

type Dispatcher struct {
	producer *kafka.Producer
}

func NewDispatcher(cfg map[string]interface{}) (events.Dispatcher, error) {
	config := &kafka.ConfigMap{}
	for k, v := range cfg {
		err := config.SetKey(k, v)
		if err != nil {
			return nil, err
		}
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, err
	}

	return &Dispatcher{producer: producer}, nil
}

func (kp *Dispatcher) Close() {
	kp.producer.Close()
}

func (kp *Dispatcher) Dispatch(ctx context.Context, topic string, msg events.Message) error {
	err := kp.producer.Produce(&kafka.Message{
		Key:            msg.Key,
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg.Body,
	}, nil)

	if err != nil {
		return err
	}

	e := <-kp.producer.Events()
	switch ev := e.(type) {
	case *kafka.Message:
		if ev.TopicPartition.Error != nil {
			log.Printf("Error delivering message: %v", ev.TopicPartition.Error)
			return ev.TopicPartition.Error
		}
	}

	return nil
}
