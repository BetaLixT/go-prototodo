package repos

import (
	"context"
	domcntxt "prototodo/pkg/domain/base/cntxt"
	infrcntxt "prototodo/pkg/infra/cntxt"
	implcntxt "prototodo/pkg/infra/impls/evcqrs/cntxt"
	"time"

	"github.com/BetaLixT/go-resiliency/retrier"
	"go.uber.org/zap"
)

type ContextFactory struct {
}

func (f *ContextFactory) Create(
	ctx context.Context,
	timeout time.Duration,
) domcntxt.IContext {
	c := &internalContext{
		done: make(chan struct{}, 1),

		rtr: *retrier.New(retrier.ExponentialBackoff(
			5,
			500*time.Millisecond,
		),
			retrier.DefaultClassifier{},
		),
		cmp: []implcntxt.Action{},
		cmt: []implcntxt.Action{},
	}
}

var _ context.Context = (*internalContext)(nil)
var _ domcntxt.IContext = (*internalContext)(nil)
var _ infrcntxt.IContext = (*internalContext)(nil)
var _ implcntxt.IContext = (*internalContext)(nil)

type internalContext struct {
	// deadline time.Time
	err  error
	done chan struct{}

	// - transaction
	rtr retrier.Retrier
	cmp []implcntxt.Action
	cmt []implcntxt.Action
}

// - Base context functions
func (c *internalContext) Deadline() (time.Time, bool) {
	return time.Now(), false
}

func (c *internalContext) Done() <-chan struct{} {
	return c.done
}

func (c *internalContext) Err() error {
	return c.err
}

func (c *internalContext) Value(key any) any {
	return nil
}

// - Transaction functions
func (c *internalContext) CommitTransaction() error {
	for _, a := range c.cmt {
		err := a()
		if err != nil {
		  return err
		}
	}
	return nil
}

func (c *internalContext) RollbackTransaction() {
	for _, a := range c.cmp {
		err := c.rtr.Run(func() error {
			err := a()
			if err != nil {
				lgr.Warn("failed to run compensatory action", zap.Error(err))
			}
			return err
		})
		if err != nil {
			lgr.Error(
				"failed to run compensatory action, max retries exceeded",
				zap.Error(err),
			)
		}
	}
}


var _ context.Context = (*minimalContext)(nil)
var _ infrcntxt.IContext = (*minimalContext)(nil)

type minimalContext struct {
	err  error
	done chan struct{}
}

// - Base context functions
func (c *minimalContext) Deadline() (time.Time, bool) {
	return time.Now(), false
}

func (c *minimalContext) Done() <-chan struct{} {
	return c.done
}

func (c *minimalContext) Err() error {
	return nil
}

func (c *minimalContext) Value(key any) any {
	return nil
}
