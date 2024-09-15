package kafka

import (
	"context"
	confluentKafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/goccy/go-json"
	"github.com/isaquesb/url-shortener/internal/events"
	"github.com/isaquesb/url-shortener/internal/ports/output"
	"github.com/isaquesb/url-shortener/pkg/logger"
)

type Dispatcher struct {
	producer *confluentKafka.Producer
}

func NewDispatcher(cfg map[string]interface{}) (output.Dispatcher, error) {
	config := &confluentKafka.ConfigMap{}
	for k, v := range cfg {
		if err := config.SetKey(k, v); err != nil {
			return nil, err
		}
	}

	producer, err := confluentKafka.NewProducer(config)
	if err != nil {
		return nil, err
	}

	return &Dispatcher{producer: producer}, nil
}

func (kp *Dispatcher) Close() {
	defer kp.producer.Close()
	kp.producer.Flush(15 * 1000)
}

func (kp *Dispatcher) Listen(_ context.Context) {
	for e := range kp.producer.Events() {
		switch ev := e.(type) {
		case *confluentKafka.Message:
			if ev.TopicPartition.Error != nil {
				logger.Error("KafkaError: delivering message", "topic", ev.TopicPartition)
			}
		}
	}
}

func (kp *Dispatcher) Dispatch(_ context.Context, msg events.Event) error {
	encoded, err := json.Marshal(events.Envelop{
		Name:  msg.GetName(),
		Event: msg,
	})
	if err != nil {
		return err
	}

	topic := msg.GetName()
	return kp.producer.Produce(&confluentKafka.Message{
		Key: msg.GetKey(),
		TopicPartition: confluentKafka.TopicPartition{
			Topic:     &topic,
			Partition: confluentKafka.PartitionAny,
		},
		Value: encoded,
	}, nil)
}
