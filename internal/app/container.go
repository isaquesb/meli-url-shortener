package app

import (
	"context"
	inputevents "github.com/isaquesb/meli-url-shortener/internal/ports/input/events"
	"github.com/isaquesb/meli-url-shortener/internal/ports/input/http"
	"github.com/isaquesb/meli-url-shortener/internal/ports/output/events"
)

type App struct {
	Ctx         context.Context
	Environment string
	Name        string
	Host        string
	Debug       Debug
	Api         *Api
	Events      map[string]string
	Worker      *Worker
}

type HasDispatcher interface {
	GetDispatcher() events.Dispatcher
}

type HasConsumer interface {
	GetConsumer() inputevents.Consumer
}

type WithDispatcher struct {
	Dispatcher       events.Dispatcher
	CreateDispatcher func() events.Dispatcher
}

type WithConsumer struct {
	GroupName      string
	Consumer       inputevents.Consumer
	CreateConsumer func() inputevents.Consumer
}

type Api struct {
	Port   int
	Router func(serviceName, environment string) http.Router
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

func (a *Api) GetDispatcher() events.Dispatcher {
	if a.WithDispatcher.Dispatcher == nil {
		a.WithDispatcher.Dispatcher = a.WithDispatcher.CreateDispatcher()
	}
	return a.WithDispatcher.Dispatcher
}

func (w *Worker) GetDispatcher() events.Dispatcher {
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
