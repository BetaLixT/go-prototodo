package repos

import (
	"context"
	domcntxt "prototodo/pkg/domain/base/cntxt"
	infrcntxt "prototodo/pkg/infra/cntxt"
	implcntxt "prototodo/pkg/infra/impls/evcqrs/cntxt"
	"time"

	"github.com/betalixt/gorr"
)

// ContextFactory to create new contexts
type ContextFactory struct{}

// NewContextFactory constructor for context factory
func NewContextFactory() *ContextFactory {
	return &ContextFactory{}
}

// Create creates a new context with timeout, transactions and trace info
func (f *ContextFactory) Create(
	string,
) domcntxt.IContext {
	c := &internalContext{}
	return c
}

var (
	_ context.Context    = (*internalContext)(nil)
	_ domcntxt.IContext  = (*internalContext)(nil)
	_ infrcntxt.IContext = (*internalContext)(nil)
	_ implcntxt.IContext = (*internalContext)(nil)
)

type internalContext struct{}

// - Base context functions
func (c *internalContext) cancel(err error) {
}

func (c *internalContext) Cancel() {
}

func (c *internalContext) Deadline() (time.Time, bool) {
	return time.Now(), false
}

func (c *internalContext) Done() <-chan struct{} {
	return nil
}

func (c *internalContext) Err() error {
	return nil
}

func (c *internalContext) Value(key any) any {
	return nil
}

// - Transaction functions
func (c *internalContext) SetTimeout(time.Duration) {
}

func (c *internalContext) CommitTransaction() error {
	return gorr.NewNotImplemented()
}

// TODO: better handling failed rollback transaction
func (c *internalContext) RollbackTransaction() {
}

func (c *internalContext) RegisterCompensatoryAction(
	cmp ...implcntxt.Action,
) {
}

func (c *internalContext) RegisterCommitAction(
	cmp ...implcntxt.Action,
) {
}

func (c *internalContext) RegisterEvent(
	id uint64,
	sagaID *string,
	stream string,
	streamID string,
	event string,
	version uint64,
	eventTime time.Time,
	data interface{},
) {
}

func (c *internalContext) GetTransactionObject(
	key string,
	constr implcntxt.Constructor,
) (interface{}, bool, error) {
	return nil, false, gorr.NewNotImplemented()
}

func (c *internalContext) GetTraceInfo() (ver, tid, pid, rid, flg string) {
	return "", "", "", "", ""
}

// - Minimal context

var (
	_ context.Context    = (*minimalContext)(nil)
	_ infrcntxt.IContext = (*minimalContext)(nil)
)

func newMinimalContext(ctx *internalContext) *minimalContext {
	return &minimalContext{}
}

type minimalContext struct{}

// - Base context functions
func (c *minimalContext) Deadline() (time.Time, bool) {
	return time.Now(), false
}

func (c *minimalContext) Done() <-chan struct{} {
	return nil
}

func (c *minimalContext) Err() error {
	return nil
}

func (c *minimalContext) Value(key any) any {
	return nil
}

func (c *minimalContext) GetTraceInfo() (ver, tid, pid, rid, flg string) {
	return "", "", "", "", ""
}
