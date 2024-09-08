package api

import (
	"context"
	"fmt"
	"github.com/isaquesb/meli-url-shortener/config"
	"github.com/isaquesb/meli-url-shortener/internal/hasher"
	router "github.com/isaquesb/meli-url-shortener/pkg/fasthttp"
	"github.com/isaquesb/meli-url-shortener/pkg/kafka"
	"github.com/isaquesb/meli-url-shortener/pkg/logger"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var appUrl string
var kp *kafka.Producer

func Start(ctx context.Context, opt *Options) {
	appUrl = config.GetEnv("APP_URL", "http://localhost:8080")
	kp, _ = kafka.NewKafkaProducer()
	defer kp.Close()

	r := router.New(
		config.GetEnv("APP_NAME", "meli-url-shortener"),
		config.Environment(),
	)

	r.POST("/", handleRequest)

	server := &fasthttp.Server{
		Handler: r.Handler,
	}

	go func(port int) {
		fmt.Printf("Starting server on :%d\n", port)
		if err := server.ListenAndServe(fmt.Sprintf(":%d", port)); err != nil {
			log.Fatalf("Error in ListenAndServe: %s", err)
		}
	}(opt.Port)

	signalListener(ctx, server)
}

func signalListener(_ context.Context, e *fasthttp.Server) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-ch

	if err := e.Shutdown(); err != nil {
		logger.Error("Error in Shutdown: %s", err)
	}
	logger.Debug("gracefully stopped http server")
}

func handleRequest(ctx *fasthttp.RequestCtx) {
	if string(ctx.Method()) != fasthttp.MethodPost {
		ctx.Error("Method not allowed", fasthttp.StatusMethodNotAllowed)
		return
	}

	url := ctx.FormValue("url")
	if len(url) == 0 {
		ctx.Error("Missing 'url' field", fasthttp.StatusBadRequest)
		return
	}

	short := hasher.GetUrlHash(string(url))
	completeUrl := fmt.Sprintf("%s/%s", appUrl, short)

	kp.Write(config.GetEnv("KAFKA_TOPIC", "shortener_urls"), string(short), string(url))

	acceptHeader := string(ctx.Request.Header.Peek("Accept"))

	if strings.Contains(acceptHeader, "application/json") {
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.SetBodyString(fmt.Sprintf(`{"short": "%s"}`, completeUrl))
		return
	}

	ctx.Response.Header.Set("Content-Type", "text/plain")
	ctx.Response.SetBodyString(completeUrl)
}
