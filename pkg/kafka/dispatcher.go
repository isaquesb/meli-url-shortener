package kafka

import (
	"context"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/goccy/go-json"
	"github.com/isaquesb/meli-url-shortener/internal/events"
	"github.com/isaquesb/meli-url-shortener/internal/ports/output"
	"log"
)

type Dispatcher struct {
	producer *ckafka.Producer
}

func NewDispatcher(cfg map[string]interface{}) (output.Dispatcher, error) {
	config := &ckafka.ConfigMap{}
	for k, v := range cfg {
		err := config.SetKey(k, v)
		if err != nil {
			return nil, err
		}
	}

	producer, err := ckafka.NewProducer(config)
	if err != nil {
		return nil, err
	}

	return &Dispatcher{producer: producer}, nil
}

func (kp *Dispatcher) Close() {
	kp.producer.Close()
}

func (kp *Dispatcher) Dispatch(_ context.Context, msg events.Event) error {
	topic := msg.GetName()
	encoded, err := json.Marshal(events.Envelop{
		Name:  msg.GetName(),
		Event: msg,
	})
	if err != nil {
		return err
	}

	err = kp.producer.Produce(&ckafka.Message{
		Key: msg.GetKey(),
		TopicPartition: ckafka.TopicPartition{
			Topic:     &topic,
			Partition: ckafka.PartitionAny,
		},
		Value: encoded,
	}, nil)

	if err != nil {
		return err
	}

	e := <-kp.producer.Events()
	switch ev := e.(type) {
	case *ckafka.Message:
		if ev.TopicPartition.Error != nil {
			log.Printf("Error delivering message: %v", ev.TopicPartition.Error)
			return ev.TopicPartition.Error
		}
	}

	return nil
}
