// Package tracectx defines functionality for parsing, injecting and extracting
// trace context
package tracectx

import "context"

// IRepository an interface that defines trace context functionality
type IRepository interface {
	ParseTraceParent(
		parent context.Context,
		traceprnt string,
	) (context.Context, error)
	ExtractTraceParent(
		ctx context.Context,
	) TxModel
}
