package repos

import (
	"context"
	"prototodo/pkg/domain/base/logger"
	"prototodo/pkg/domain/base/tracectx"

	"github.com/betalixt/gorr"
)

var _ tracectx.IRepository = (*TraceRepository)(nil)

// TraceRepository repository for generating, injecting and extracting trace
// info
type TraceRepository struct {
	lgrf logger.IFactory
}

// NewTraceRepository constructs a TraceRepository
func NewTraceRepository() *TraceRepository {
	return &TraceRepository{}
}

// ParseTraceParent parses and or generates trace information and returns
// context with trace information injected
func (r *TraceRepository) ParseTraceParent(
	parent context.Context,
	traceprnt string,
) (context.Context, error) {
	return nil, gorr.NewNotImplemented()
}

// ExtractTraceParent extracts injected trace information from context
func (*TraceRepository) ExtractTraceParent(
	ctx context.Context,
) tracectx.TxModel {
	return tracectx.TxModel{}
}
