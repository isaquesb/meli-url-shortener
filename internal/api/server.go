package api

import (
	"context"
	"fmt"
	"github.com/isaquesb/meli-url-shortener/internal/ports/input/http"
	"github.com/isaquesb/meli-url-shortener/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Start(_ context.Context, server http.Server, router http.Router) {
	router.POST("/", CreateShortUrl)

	go func() {
		logger.Info(fmt.Sprintf("starting http server on port %d", server.Options().Port))
		if err := server.Start(router); err != nil {
			log.Fatalf("error starting server: %s", err)
		}
	}()

	signalListener(server)
}

func signalListener(e http.Server) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-ch

	if err := e.Shutdown(); err != nil {
		logger.Error("Error in Shutdown: %s", err)
	}
	logger.Info("gracefully stopped http server")
}
