package api

import (
	"context"
	"fmt"
	"github.com/isaquesb/url-shortener/internal/ports/input/http"
	"github.com/isaquesb/url-shortener/internal/urls"
	"github.com/isaquesb/url-shortener/pkg/instrumentation"
	"github.com/isaquesb/url-shortener/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	Ctx    context.Context
	Http   http.Server
	Router http.Router
	Instr  *instrumentation.Instrumentation
}

func (a *Server) Start() {
	a.Router.GET("/{short}", urls.RedirectShort)
	a.Router.DELETE("/{short}", urls.DeleteShortUrl)
	a.Router.GET("/{short}/stats", urls.ShowStats)
	a.Router.POST("/", urls.CreateShortUrl)

	go func() {
		logger.Info(fmt.Sprintf("starting http server on port %d", a.Http.Options().Port))
		if err := a.Http.Start(a.Router); err != nil {
			log.Fatalf("error starting server: %s", err)
		}
	}()

	a.signalListener()
}

func (a *Server) signalListener() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-ch

	if err := a.Instr.Tracer.Shutdown(a.Ctx); err != nil {
		logger.Error("Failed to shutdown tracer: %v", err)
	}

	if err := a.Http.Shutdown(); err != nil {
		logger.Error("Error in Shutdown: %s", err)
	}

	logger.Info("gracefully stopped http server")
}
