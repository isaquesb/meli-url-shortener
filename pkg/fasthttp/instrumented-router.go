package fasthttp

import (
	"github.com/isaquesb/meli-url-shortener/internal/ports/input/http"
	"github.com/isaquesb/meli-url-shortener/pkg/instrumentation"
	"github.com/isaquesb/meli-url-shortener/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
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

func NewRouter(instrumentation *instrumentation.Instrumentation) *Router {
	return &Router{
		router:          fr.New(),
		instrumentation: instrumentation,
	}
}

func (ir *Router) instrumentedHandler(handlerFunc fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()

		tracer := otel.Tracer("fast_http-tracer")
		_, span := tracer.Start(ctx, "httpHandler", trace.WithAttributes(
			attribute.String("method", string(ctx.Method())),
			attribute.String("path", string(ctx.Path())),
		))
		defer span.End()

		handlerFunc(ctx)
		duration := time.Since(start).Milliseconds()

		labels := []attribute.KeyValue{
			attribute.String("method", string(ctx.Method())),
			attribute.Int("status", ctx.Response.StatusCode()),
			attribute.String("path", string(ctx.Path())),
		}

		ir.instrumentation.Metrics.HTTPRequestDuration.Record(ctx, duration, api.WithAttributes(labels...))
		ir.instrumentation.Metrics.HTTPTotalRequestsCounter.Add(ctx, 1, api.WithAttributes(labels...))
	}
}

func (ir *Router) routedHandler(handlerFunc http.RouteHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		httpRequest := NewRequest(ctx)
		response, err := handlerFunc(httpRequest)
		if err != nil {
			logger.Error(err.Error())
			ctx.Error("internal server error", fasthttp.StatusInternalServerError)
			return
		}
		ParseResponse(response, ctx)
	}
}

func (ir *Router) Handler(ctx *fasthttp.RequestCtx) {
	ir.router.Handler(ctx)
}

func (ir *Router) GET(path string, handler http.RouteHandler) {
	ir.router.GET(path, ir.instrumentedHandler(ir.routedHandler(handler)))
}

func (ir *Router) POST(path string, handler http.RouteHandler) {
	ir.router.POST(path, ir.instrumentedHandler(ir.routedHandler(handler)))
}

func (ir *Router) DELETE(path string, handler http.RouteHandler) {
	ir.router.DELETE(path, ir.instrumentedHandler(ir.routedHandler(handler)))
}
