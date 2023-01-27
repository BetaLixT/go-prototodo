package logger

import (
	"context"
	"go.uber.org/zap"
)

type LoggerFactory struct {
	lgr *zap.Logger
}

func NewLoggerFactory() (*LoggerFactory, error) {
	lgr, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &LoggerFactory{
		lgr: lgr,
	}, nil
}

func (lf *LoggerFactory) NewLogger(
	ctx context.Context,
) *zap.Logger {
	if ctx == nil {
		return lf.lgr
	}
	raw := ctx.Value(common.TRACE_INFO_KEY)

	if raw == nil {
		return lf.lgr
	}
	trace, ok := raw.(trex.TxModel)
	if !ok {
		return lf.lgr
	}
	return lf.lgr.With(
		zap.String("tid", trace.Tid),
		zap.String("pid", trace.Pid),
		zap.String("rid", trace.Rid),
	)
}

func (lf *LoggerFactory) Close() {
	lf.lgr.Sync()
}
