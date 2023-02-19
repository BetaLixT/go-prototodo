package appinsights

import (
	"github.com/Soreing/apex"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type AppInsightsTraceExporter sdktrace.SpanExporter

func NewAppInsightsTraceExporter(
	opts *ExporterOptions,
) (AppInsightsTraceExporter, error) {
	if opts.InstrKey == "" {
		return nil, nil
	}

	return apex.NewExporter(
		opts.InstrKey,
		func(msg string) error {
			return nil
		},
	)
}
