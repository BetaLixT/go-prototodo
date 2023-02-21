// Package promex adds prometheus metrix
package promex

import (
	"context"
	"prototodo/pkg/domain/common"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// TraceExporter an application insights exporter
type TraceExporter sdktrace.SpanExporter

// NewTraceExporter constructs an app insights exporter
func NewTraceExporter() (TraceExporter, error) {
	return NewExporter(common.ServiceName), nil
}

type Exporter struct {
	requests      prometheus.Counter
	requestStatus prometheus.CounterVec
	responseTime  prometheus.HistogramVec

	events        prometheus.Counter
	eventsSuccess prometheus.CounterVec
	eventTime     prometheus.HistogramVec

	depStatus prometheus.CounterVec
}

func NewExporter(prefix string) *Exporter {
	return &Exporter{
		requests: promauto.NewCounter(prometheus.CounterOpts{
			Name: prefix + "_processed_reqs_total",
			Help: "The total number of processed requests",
		}),
		requestStatus: *promauto.NewCounterVec(prometheus.CounterOpts{
			Name: prefix + "_processed_reqs_status",
			Help: "The status codes of requests",
		}, []string{"code"}),
		responseTime: *promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name: prefix + "_processed_reqs_latency",
			Help: "The latency of requests",
		}, []string{"uri"}),

		events: promauto.NewCounter(prometheus.CounterOpts{
			Name: prefix + "_processed_evnts_total",
			Help: "The total number of processed events",
		}),
		eventsSuccess: *promauto.NewCounterVec(prometheus.CounterOpts{
			Name: prefix + "_processed_evnts_status",
			Help: "The status codes of events",
		}, []string{"status"}),
		eventTime: *promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name: prefix + "_processed_evnts_latency",
			Help: "The latency of events",
		}, []string{"key"}),

		depStatus: *promauto.NewCounterVec(prometheus.CounterOpts{
			Name: prefix + "_dependencies_status",
			Help: "The status of dependencies",
		}, []string{"status", "type"}),
	}
}

func (exp *Exporter) Shutdown(ctx context.Context) error {
	return nil
}

func (exp *Exporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	for i := range spans {
		exp.process(spans[i])
	}
	return nil
}

// Preprocesses the Otel span and dispatches it to app insights differently
// based on the span kind.
func (exp *Exporter) process(sp sdktrace.ReadOnlySpan) {
	success := true
	if sp.Status().Code != codes.Ok {
		success = false
	}

	props := map[string]string{}

	rattr := sp.Resource().Attributes()
	for _, e := range rattr {
		props[string(e.Key)] = e.Value.AsString()
	}
	attr := sp.Attributes()
	for _, e := range attr {
		props[string(e.Key)] = e.Value.AsString()
	}

	switch sp.SpanKind() {
	case trace.SpanKindServer:
		exp.processRequest(sp, success, props)
	case trace.SpanKindClient:
		exp.processDependency(sp, success, props)
	case trace.SpanKindProducer:
		exp.processDependency(sp, success, props)
	case trace.SpanKindConsumer:
		exp.processEvent(sp, success, props)
	}
}

func (exp *Exporter) processRequest(
	sp sdktrace.ReadOnlySpan,
	success bool,
	properties map[string]string,
) {
	exp.requests.Inc()

	if val, ok := properties["url"]; ok {
		delete(properties, "url")
		exp.responseTime.WithLabelValues(val).
			Observe(sp.EndTime().Sub(sp.StartTime()).Seconds())
	}
	if val, ok := properties["responseCode"]; ok {
		delete(properties, "responseCode")
		exp.requestStatus.WithLabelValues(val).Inc()
	}
}

func (exp *Exporter) processEvent(
	sp sdktrace.ReadOnlySpan,
	success bool,
	properties map[string]string,
) {
	exp.events.Inc()

	if val, ok := properties["key"]; ok {
		delete(properties, "key")
		exp.eventTime.WithLabelValues(val).
			Observe(sp.EndTime().Sub(sp.StartTime()).Seconds())
	}
	if success {
		exp.eventsSuccess.WithLabelValues("success").Inc()
	} else {
		exp.eventsSuccess.WithLabelValues("failed").Inc()
	}
}

func (exp *Exporter) processDependency(
	sp sdktrace.ReadOnlySpan,
	success bool,
	properties map[string]string,
) {
	typ := ""
	if val, ok := properties["type"]; ok {
		delete(properties, "type")
		typ = val
	}
	if success {
		exp.depStatus.WithLabelValues("success", typ).Inc()
	} else {
		exp.depStatus.WithLabelValues("failed", typ).Inc()
	}
}
