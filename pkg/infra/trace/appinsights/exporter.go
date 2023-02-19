package appinsights

import (
	"techunicorn.com/udc-core/pretz/pkg/infra/logger"

	"github.com/Soreing/apex"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type AppInsightsTraceExporter sdktrace.SpanExporter

func NewAppInsightsTraceExporter(
	opts *ExporterOptions,
	lgrf *logger.LoggerFactory,
) (AppInsightsTraceExporter, error) {
	if opts.InstrKey == "" {
		return nil, nil
	}

	return apex.NewExporter(
		opts.InstrKey,
		opts.ServiceName,
		func(msg string) error {
			return nil
		},
	)
}
