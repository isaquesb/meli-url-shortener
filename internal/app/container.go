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
	Topics          map[string]string
	Worker          *Worker
	Instrumentation func() *instrumentation.Instrumentation
}

type Lazy[T any] struct {
	Instance T
	Create   func() T
}

func (l *Lazy[T]) Get() T {
	if l.isNil() {
		l.Instance = l.Create()
	}
	return l.Instance
}

func (l *Lazy[T]) isNil() bool {
	var t T
	return any(l.Instance) == any(t)
}

type Api struct {
	Port       int
	Router     func(*instrumentation.Instrumentation) http.Router
	Server     func(options http.Options) http.Server
	Dispatcher Lazy[output.Dispatcher]
	Repository Lazy[output.UrlRepository]
}

type Worker struct {
	Dispatcher Lazy[output.Dispatcher]
	Consumer   Lazy[inputevents.Consumer]
}

type Debug struct {
	Enabled bool
	Trace   bool
}

func GetApp() *App {
	return container
}

func SetApp(app *App) {
	container = app
}
