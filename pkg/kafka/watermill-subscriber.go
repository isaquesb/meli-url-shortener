package kafka

import (
	"github.com/ThreeDotsLabs/watermill"
	watermillKafka "github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewSubscriber(logger watermill.LoggerAdapter, consumerGroup string, brokers []string) (message.Subscriber, error) {
	sub, err := watermillKafka.NewSubscriber(watermillKafka.SubscriberConfig{
		Brokers:       brokers,
		Unmarshaler:   watermillKafka.DefaultMarshaler{},
		ConsumerGroup: consumerGroup,
	}, logger)

	if err != nil {
		return nil, err
	}
	return sub, nil
}
