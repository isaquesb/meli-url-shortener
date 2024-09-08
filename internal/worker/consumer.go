package worker

import (
	"context"
	"github.com/isaquesb/meli-url-shortener/config"
)

func Consume(ctx context.Context) {
	wLogger := config.MessageRouterLogger()
	router := config.MessageRouter(wLogger)
	sub, err := config.KafkaSubscriber(wLogger)

	if err != nil {
		panic(err)
	}

	router.AddNoPublisherHandler(
		"store_database",
		config.GetEnv("KAFKA_TOPIC", "shortener_urls"),
		sub,
		StoreDatabase,
	)

	err = router.Run(ctx)
	if err != nil {
		panic(err)
	}
}
