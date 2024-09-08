package instrumentation

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

type Instrumentation struct {
	HTTPTotalRequestsCounter api.Int64Counter
	HTTPRequestHistogram     api.Float64Histogram
}

func New(serviceName, environment string) *Instrumentation {
	exporter, err := otlpmetrichttp.New(context.Background(), otlpmetrichttp.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to create the collector exporter: %v", err)
	}

	reader := metric.NewPeriodicReader(exporter, metric.WithInterval(15*time.Second))

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		attribute.String("environment", environment),
	)

	// Initialize the meter provider
	provider := metric.NewMeterProvider(
		metric.WithReader(reader),
		metric.WithResource(res),
	)

	meter := provider.Meter("fasthttp-otel")

	// Create instruments for our metrics
	httpTotalRequestsCounter, _ := meter.Int64Counter(
		"http_total_requests",
		api.WithDescription("Total number of HTTP requests"),
	)

	httpRequestHistogram, _ := meter.Float64Histogram(
		"http_request_duration_seconds",
		api.WithDescription("Duration of HTTP requests"),
	)

	return &Instrumentation{
		HTTPTotalRequestsCounter: httpTotalRequestsCounter,
		HTTPRequestHistogram:     httpRequestHistogram,
	}
}
