package trace

import "context"

type IRepository interface {
	// Parses (or generates new) traceparent and returns
	// context with trace info injected
	ParseTraceParent(
		parent context.Context,
		traceprnt string,
	) (context.Context, error)
	ExtractTraceParent(
		ctx context.Context,
	) TxModel
}
