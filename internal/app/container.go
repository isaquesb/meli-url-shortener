package app

import (
	"context"
	inputevents "github.com/isaquesb/meli-url-shortener/internal/ports/input/events"
	"github.com/isaquesb/meli-url-shortener/internal/ports/input/http"
	"github.com/isaquesb/meli-url-shortener/internal/ports/output"
	"github.com/isaquesb/meli-url-shortener/pkg/instrumentation"
)

var container *App

type App struct {
	Ctx             context.Context
	Environment     string
	Name            string
	Host            string
	Debug           Debug
	Api             *Api
	Events          map[string]string
	Worker          *Worker
	Instrumentation func() *instrumentation.Instrumentation
}

type HasDispatcher interface {
	GetDispatcher() output.Dispatcher
}

type HasConsumer interface {
	GetConsumer() inputevents.Consumer
}

type WithDispatcher struct {
	Dispatcher       output.Dispatcher
	CreateDispatcher func() output.Dispatcher
}

type WithConsumer struct {
	GroupName      string
	Consumer       inputevents.Consumer
	CreateConsumer func() inputevents.Consumer
}

type Api struct {
	Port   int
	Router func(*instrumentation.Instrumentation) http.Router
	Server func(options http.Options) http.Server
	WithDispatcher
}

type Worker struct {
	Consumer inputevents.Consumer
	WithDispatcher
	WithConsumer
}

type Debug struct {
	Enabled bool
	Trace   bool
}

func (a *Api) GetDispatcher() output.Dispatcher {
	if a.WithDispatcher.Dispatcher == nil {
		a.WithDispatcher.Dispatcher = a.WithDispatcher.CreateDispatcher()
	}
	return a.WithDispatcher.Dispatcher
}

func (w *Worker) GetDispatcher() output.Dispatcher {
	if w.WithDispatcher.Dispatcher == nil {
		w.WithDispatcher.Dispatcher = w.WithDispatcher.CreateDispatcher()
	}
	return w.WithDispatcher.Dispatcher
}

func (w *Worker) GetConsumer() inputevents.Consumer {
	if w.WithConsumer.Consumer == nil {
		w.WithConsumer.Consumer = w.WithConsumer.CreateConsumer()
	}
	return w.WithConsumer.Consumer
}

func GetApp() *App {
	return container
}

func SetApp(app *App) {
	container = app
}
