package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

type Producer struct {
	producer *kafka.Producer
}

func NewKafkaProducer() (*Producer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, err
	}

	return &Producer{producer: producer}, nil
}

func (kp *Producer) Close() {
	kp.producer.Close()
}

func (kp *Producer) Write(topic, short, url string) error {
	message := fmt.Sprintf(`%s%s`, short, url)

	err := kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
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
