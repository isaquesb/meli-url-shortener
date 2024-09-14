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

func Start(ctx context.Context, server http.Server, router http.Router, instrumentation *instrumentation.Instrumentation) {
	router.POST("/", urls.CreateShortUrl)

	go func() {
		logger.Info(fmt.Sprintf("starting http server on port %d", server.Options().Port))
		if err := server.Start(router); err != nil {
			log.Fatalf("error starting server: %s", err)
		}
	}()

	signalListener(ctx, server, instrumentation)
}

func signalListener(ctx context.Context, server http.Server, instrumentation *instrumentation.Instrumentation) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-ch

	if err := instrumentation.Tracer.Shutdown(ctx); err != nil {
		logger.Error("Failed to shutdown tracer: %v", err)
	}

	if err := server.Shutdown(); err != nil {
		logger.Error("Error in Shutdown: %s", err)
	}

	logger.Info("gracefully stopped http server")
}
