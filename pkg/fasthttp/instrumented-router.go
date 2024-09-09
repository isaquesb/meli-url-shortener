package fasthttp

import (
	"github.com/isaquesb/meli-url-shortener/pkg/instrumentation"
	"time"

	fr "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/metric"
)

type Router struct {
	router          *fr.Router
	instrumentation *instrumentation.Instrumentation
}

func New(serviceName, environment string) *Router {
	return &Router{
		router:          fr.New(),
		instrumentation: instrumentation.New(serviceName, environment),
	}
}

func (ir *Router) instrumentedHandler(handlerFunc fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		handlerFunc(ctx)
		duration := time.Since(start).Milliseconds()

		labels := []attribute.KeyValue{
			attribute.String("method", string(ctx.Method())),
			attribute.Int("status", ctx.Response.StatusCode()),
			attribute.String("path", string(ctx.Path())),
		}

		ir.instrumentation.HTTPRequestDuration.Record(ctx, duration, api.WithAttributes(labels...))
		ir.instrumentation.HTTPTotalRequestsCounter.Add(ctx, 1, api.WithAttributes(labels...))
	}
}

func (ir *Router) Handler(ctx *fasthttp.RequestCtx) {
	ir.router.Handler(ctx)
}

func (ir *Router) GET(path string, handler fasthttp.RequestHandler) {
	ir.router.GET(path, ir.instrumentedHandler(handler))
}

func (ir *Router) POST(path string, handler fasthttp.RequestHandler) {
	ir.router.POST(path, ir.instrumentedHandler(handler))
}

func (ir *Router) DELETE(path string, handler fasthttp.RequestHandler) {
	ir.router.DELETE(path, ir.instrumentedHandler(handler))
}
