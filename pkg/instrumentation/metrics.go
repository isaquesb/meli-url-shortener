package instrumentation

import (
	"context"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	metricApi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"log"
	"time"
)

type Metrics struct {
	HTTPTotalRequestsCounter metricApi.Int64Counter
	HTTPRequestDuration      metricApi.Int64Gauge
}

func NewMetrics(ctx context.Context, res *resource.Resource) *Metrics {
	exporter, err := otlpmetrichttp.New(ctx, otlpmetrichttp.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to create the collector exporter: %v", err)
	}

	reader := metric.NewPeriodicReader(exporter, metric.WithInterval(15*time.Second))

	provider := metric.NewMeterProvider(
		metric.WithReader(reader),
		metric.WithResource(res),
	)

	meter := provider.Meter("fasthttp-otel")

	// Create instruments for our metrics
	httpTotalRequestsCounter, _ := meter.Int64Counter(
		"http_total_requests",
		metricApi.WithDescription("Total number of HTTP requests"),
	)

	httpRequestDuration, _ := meter.Int64Gauge(
		"http_request_duration_ms",
		metricApi.WithDescription("Duration of HTTP requests"),
	)

	return &Metrics{
		HTTPTotalRequestsCounter: httpTotalRequestsCounter,
		HTTPRequestDuration:      httpRequestDuration,
	}
}
