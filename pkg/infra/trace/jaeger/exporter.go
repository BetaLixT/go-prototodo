package jaeger

import (
	jexp "go.opentelemetry.io/otel/exporters/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type JaegerTraceExporter sdktrace.SpanExporter

func NewJaegerTraceExporter(
	opts *ExporterOptions,
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
