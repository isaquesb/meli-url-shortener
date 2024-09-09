package config

import (
	"context"
	"github.com/isaquesb/meli-url-shortener/internal/app"
	inputevents "github.com/isaquesb/meli-url-shortener/internal/ports/input/events"
	"github.com/isaquesb/meli-url-shortener/internal/ports/input/http"
	"github.com/isaquesb/meli-url-shortener/internal/ports/output/events"
	"github.com/isaquesb/meli-url-shortener/pkg/dynamoDb"
	"github.com/isaquesb/meli-url-shortener/pkg/fasthttp"
	"github.com/isaquesb/meli-url-shortener/pkg/kafka"
)

var container *app.App

func NewApp() *app.App {
	return &app.App{
		Ctx:         context.TODO(),
		Environment: Environment(),
		Name:        GetEnv("APP_NAME", "meli-url-shortener"),
		Host:        GetEnv("APP_HOST", "http://localhost:8080"),
		Events: map[string]string{
			"urls.created": GetEnv("CREATED_TOPIC", "urls.created"),
			"urls.deleted": GetEnv("DELETED_TOPIC", "urls.deleted"),
			"urls.visited": GetEnv("VISITED_TOPIC", "urls.visited"),
		},
		Debug: app.Debug{
			Enabled: GetEnv("APP_DEBUG", "false") == "true",
			Trace:   GetEnv("APP_TRACE", "false") == "true",
		},
		Api: &app.Api{
			Port: GetIntEnv("API_PORT", "8080"),
			Router: func(serviceName, environment string) http.Router {
				return fasthttp.NewRouter(serviceName, environment)
			},
			Server: func(options http.Options) http.Server {
				return fasthttp.New(options)
			},
			WithDispatcher: app.WithDispatcher{
				Dispatcher: nil,
				CreateDispatcher: func() events.Dispatcher {
					kafkaDispatcher, err := kafka.NewDispatcher(map[string]interface{}{
						"bootstrap.servers": GetEnv("KAFKA_BROKERS", "kafka:9092"),
					})
					if err != nil {
						panic(err)
					}
					return kafkaDispatcher
				},
			},
		},
		Worker: &app.Worker{
			WithConsumer: app.WithConsumer{
				Consumer: nil,
				CreateConsumer: func() inputevents.Consumer {
					cLogger := kafka.NewLogger(container.Debug.Enabled, container.Debug.Trace)
					group := GetEnv("KAFKA_CONSUMER_GROUP", "url-shortener")
					subscriber, err := kafka.NewSubscriber(
						cLogger,
						group,
						[]string{GetEnv("KAFKA_BROKERS", "kafka:9092")},
					)
					if err != nil {
						panic(err)
					}
					return kafka.NewConsumer(cLogger, subscriber)
				},
			},
			WithDispatcher: app.WithDispatcher{
				Dispatcher: nil,
				CreateDispatcher: func() events.Dispatcher {
					dispatcher, err := dynamoDb.NewDispatcher(dynamoDb.DispatcherOptions{
						Region: GetEnv("DYNAMODB_REGION", "local"),
						Host:   GetEnv("DYNAMODB_HOST", "http://dynamodb-local:8000"),
					})
					if err != nil {
						panic(err)
					}
					return dispatcher
				},
			},
		},
	}
}

func GetApp() *app.App {
	if container == nil {
		container = NewApp()
	}
	return container
}

func SetApp(app *app.App) {
	container = app
}
