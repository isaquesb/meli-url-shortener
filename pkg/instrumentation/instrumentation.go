package instrumentation

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

type Instrumentation struct {
	Metrics *Metrics
	Tracer  *trace.TracerProvider
}

func New(ctx context.Context, serviceName, environment string) *Instrumentation {
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		attribute.String("environment", environment),
	)

	return &Instrumentation{
		Metrics: NewMetrics(ctx, res),
		Tracer:  NewTracer(ctx, res),
	}
}
