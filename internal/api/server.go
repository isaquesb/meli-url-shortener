package api

import (
	"context"
	"fmt"
	"github.com/isaquesb/meli-url-shortener/internal/ports/input/http"
	"github.com/isaquesb/meli-url-shortener/internal/urls"
	"github.com/isaquesb/meli-url-shortener/pkg/instrumentation"
	"github.com/isaquesb/meli-url-shortener/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Api struct {
	Ctx    context.Context
	Server http.Server
	Router http.Router
	Instr  *instrumentation.Instrumentation
}

func (a *Api) Start() {
	a.Router.GET("/{short}", urls.RedirectShort)
	a.Router.DELETE("/{short}", urls.DeleteShortUrl)
	a.Router.GET("/{short}/stats", urls.ShowStats)
	a.Router.POST("/", urls.CreateShortUrl)

	go func() {
		logger.Info(fmt.Sprintf("starting http server on port %d", a.Server.Options().Port))
		if err := a.Server.Start(a.Router); err != nil {
			log.Fatalf("error starting server: %s", err)
		}
	}()

	a.signalListener()
}

func (a *Api) signalListener() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-ch

	if err := a.Instr.Tracer.Shutdown(a.Ctx); err != nil {
		logger.Error("Failed to shutdown tracer: %v", err)
	}

	if err := a.Server.Shutdown(); err != nil {
		logger.Error("Error in Shutdown: %s", err)
	}

	logger.Info("gracefully stopped http server")
}
