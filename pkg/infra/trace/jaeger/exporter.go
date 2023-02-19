package jaeger

import (
	"techunicorn.com/udc-core/pretz/pkg/infra/logger"

	jexp "go.opentelemetry.io/otel/exporters/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type JaegerTraceExporter sdktrace.SpanExporter

func NewJaegerTraceExporter(
	opts *ExporterOptions,
	lgrf *logger.LoggerFactory,
) (JaegerTraceExporter, error) {
	if opts.Endpoint == "" {
		return nil, nil
	}
	return jexp.New(
		jexp.WithCollectorEndpoint(
			jexp.WithEndpoint(opts.Endpoint),
		),
	)
}
