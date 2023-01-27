package repos

import (
	"context"
	"fmt"
	domcntxt "prototodo/pkg/domain/base/cntxt"
	"prototodo/pkg/domain/base/logger"
	infrcntxt "prototodo/pkg/infra/cntxt"
	implcntxt "prototodo/pkg/infra/impls/evcqrs/cntxt"
	"sync"
	"time"

	"github.com/BetaLixT/go-resiliency/retrier"
	"go.uber.org/zap"
)

type ContextFactory struct {
	lgrf logger.IFactory
}

func (f *ContextFactory) Create(
	ctx context.Context,
	timeout time.Duration,
) (domcntxt.IContext, context.CancelFunc) {
	c := &internalContext{
		lgrf: f.lgrf,
		err: nil,
		done: make(chan struct{}, 1),

		rtr: *retrier.New(retrier.ExponentialBackoff(
			5,
			500*time.Millisecond,
		),
			retrier.DefaultClassifier{},
		),
		compensatoryActions: []implcntxt.Action{},
		commitActions: []implcntxt.Action{},
		events: []dispatchableEvent{},
		isCommited: false,
		isRolledback: false,
		commitmtx: &sync.Mutex{},
	}

	// TODO: tracing values

	ctx, cancel := context.WithTimeout(
	  c,
	  timeout,
	)
	
}

var _ context.Context = (*internalContext)(nil)
var _ domcntxt.IContext = (*internalContext)(nil)
var _ infrcntxt.IContext = (*internalContext)(nil)
var _ implcntxt.IContext = (*internalContext)(nil)

type internalContext struct {
	lgrf logger.IFactory
	// deadline time.Time
	err  error
	done chan struct{}

	// - transaction
	rtr                 retrier.Retrier
	compensatoryActions []implcntxt.Action
	commitActions       []implcntxt.Action
	events              []dispatchableEvent
	isCommited          bool
	isRolledback        bool
	commitmtx           *sync.Mutex
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
	c.commitmtx.Lock()
	defer c.commitmtx.Unlock()
	if c.isCommited || c.isRolledback {
	  return fmt.Errorf(
	    "tried to commit transaction that has already been commited/rolled back",
	  )
	}
	ctx := newMinimalContext(c)
	for range c.events {
	  // TODO Event handling
		// err := c.ndisp.DispatchEventNotification(
		// 	ctx,
		// 	evnt.stream,
		// 	evnt.streamId,
		// 	evnt.event,
		// 	evnt.version,
		// 	evnt.data,
		// 	evnt.eventTime,
		// )
		// if err != nil {
		// 	return err
		// }
	}
	for _, commit := range c.commitActions {
		err := commit(ctx)
		if err != nil {
			return err
		}
	}
	// TODO: confirm event
	c.isCommited = true
	return nil
}

// TODO: better handling failed rollback transaction
func (c *internalContext) RollbackTransaction() {
  c.commitmtx.Lock()
	defer c.commitmtx.Unlock()
	if c.isCommited || c.isRolledback {
	  return
	}


  c.isRolledback = true
	ctx := newMinimalContext(c)
	lgr := c.lgrf.Create(ctx)
	for _, cmp := range c.compensatoryActions {
		err := c.rtr.Run(func() error {
			err := cmp(ctx)
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

func (c *internalContext) RegisterCompensatoryAction(
	cmp ...implcntxt.Action,
) {
	c.compensatoryActions = append(c.compensatoryActions, cmp...)
}

func (c *internalContext) RegisterCommitAction(
	cmp ...implcntxt.Action,
) {
	c.commitActions = append(c.commitActions, cmp...)
}

func (c *internalContext) RegisterEvent(
	id uint64,
	sagaId *string,
	stream string,
	streamId string,
	event string,
	version uint64,
	eventTime time.Time,
	data interface{},
) {
	c.events = append(c.events, dispatchableEvent{
		stream:    stream,
		streamId:  streamId,
		event:     event,
		version:   int(version),
		eventTime: eventTime,
		data:      data,
	})
}

func (c *internalContext) GetTraceInfo() (ver, tid, pid, rid, flg string) {
	return c.ver, c.tid, c.pid, c.rid, c.flg
}

// - Minimal context

var _ context.Context = (*minimalContext)(nil)
var _ infrcntxt.IContext = (*minimalContext)(nil)

func newMinimalContext(ctx *internalContext) *minimalContext {
	return &minimalContext{
		done: make(chan struct{}, 1),
		ver:  ctx.ver,
		tid:  ctx.tid,
		pid:  ctx.pid,
		rid:  ctx.rid,
		flg:  ctx.flg,
	}
}

type minimalContext struct {
	done chan struct{}
	ver  string
	tid  string
	pid  string
	rid  string
	flg  string
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

func (c *minimalContext) GetTraceInfo() (ver, tid, pid, rid, flg string) {
	return c.ver, c.tid, c.pid, c.rid, c.flg
}

type dispatchableEvent struct {
	stream    string
	streamId  string
	version   int
	event     string
	eventTime time.Time
	data      interface{}
}
