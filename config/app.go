package config

import (
	"context"
	"github.com/isaquesb/url-shortener/internal/app"
	inputevents "github.com/isaquesb/url-shortener/internal/ports/input/events"
	"github.com/isaquesb/url-shortener/internal/ports/input/http"
	"github.com/isaquesb/url-shortener/internal/ports/output"
	"github.com/isaquesb/url-shortener/internal/urls"
	"github.com/isaquesb/url-shortener/pkg/dynamoDb"
	"github.com/isaquesb/url-shortener/pkg/fasthttp"
	"github.com/isaquesb/url-shortener/pkg/instrumentation"
	"github.com/isaquesb/url-shortener/pkg/watermill"
)

func NewApp() *app.App {
	instance := &app.App{
		Ctx:         context.TODO(),
		Environment: Environment(),
		Name:        GetEnv("APP_NAME", "url-shortener"),
		Host:        GetEnv("APP_HOST", "http://localhost:8080"),
		Instrumentation: func() *instrumentation.Instrumentation {
			container := app.GetApp()
			return instrumentation.New(container.Ctx, container.Name, container.Environment)
		},
		Topics: map[string]string{
			urls.Created: GetEnv("CREATED_STORE", urls.Created),
			urls.Deleted: GetEnv("DELETED_STORE", urls.Deleted),
			urls.Visited: GetEnv("VISITED_STORE", urls.Visited),
		},
		Debug: app.Debug{
			Enabled: GetEnv("APP_DEBUG", "false") == "true",
			Trace:   GetEnv("APP_TRACE", "false") == "true",
		},
		Api: &app.Api{
			Port: GetIntEnv("API_PORT", "8080"),
			Router: func(instrumentation *instrumentation.Instrumentation) http.Router {
				return fasthttp.NewRouter(instrumentation)
			},
			Server: func(options http.Options) http.Server {
				return fasthttp.New(options)
			},
			Repository: app.Lazy[output.UrlRepository]{
				Create: func() output.UrlRepository {
					client, err := dynamoDbClient()
					if err != nil {
						panic(err)
					}
					return dynamoDb.NewRepository(client)
				},
			},
			Dispatcher: app.Lazy[output.Dispatcher]{
				Create: func() output.Dispatcher {
					dispatcher, err := GetAppDispatcher()
					if err != nil {
						panic(err)
					}
					return dispatcher
				},
			},
		},
		Worker: &app.Worker{
			Consumer: app.Lazy[inputevents.Consumer]{
				Create: func() inputevents.Consumer {
					container := app.GetApp()
					logger := watermill.NewLogger(container.Debug.Enabled, container.Debug.Trace)
					group := GetEnv("CONSUMER_GROUP", "url-shortener")
					consumer, err := GetAppConsumer(logger, group)
					if err != nil {
						panic(err)
					}
					return consumer
				},
			},
			Dispatcher: app.Lazy[output.Dispatcher]{
				Create: func() output.Dispatcher {
					client, err := dynamoDbClient()
					if err != nil {
						panic(err)
					}
					dispatcher, err := dynamoDb.NewDispatcher(client)
					if err != nil {
						panic(err)
					}
					return dispatcher
				},
			},
		},
	}
	app.SetApp(instance)

	return instance
}
