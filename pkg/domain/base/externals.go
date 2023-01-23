package base

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type IContextFactory interface {
	Create(
		ctx context.Context,
		timeout time.Duration,
	) IContext
}

type ILoggerFactory interface {
	Create(ctx context.Context) zap.Logger
}

type IContext interface {
	context.Context
}
