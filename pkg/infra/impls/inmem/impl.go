// Package inmem Event source CQRS implementation of the domain layer
package inmem

import (
	"context"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/acl"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/cntxt"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/foreigns"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/impl"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/logger"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/trace"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/uids"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/uniques"
	"techunicorn.com/udc-core/prototodo/pkg/domain/domains/quotes"
	"techunicorn.com/udc-core/prototodo/pkg/domain/domains/tasks"
	"techunicorn.com/udc-core/prototodo/pkg/infra/config"
	"techunicorn.com/udc-core/prototodo/pkg/infra/impls/inmem/repos"
	"techunicorn.com/udc-core/prototodo/pkg/infra/lgr"
	"techunicorn.com/udc-core/prototodo/pkg/infra/sf"
	"time"

	"github.com/betalixt/gorr"
	"github.com/google/wire"
)

// DependencySet dependencies provided by the implementation
var DependencySet = wire.NewSet(
	NewImplementation,
	wire.Bind(
		new(impl.IImplementation),
		new(*Implementation),
	),

	// Infra
	config.NewInitializer,
	lgr.NewLoggerFactory,
	wire.Bind(
		new(logger.IFactory),
		new(*lgr.LoggerFactory),
	),
	sf.NewSnowflake,
	config.NewSnowflakeOptions,

	// Repos
	repos.NewACLRepository,
	wire.Bind(
		new(acl.IRepository),
		new(*repos.ACLRepository),
	),
	repos.NewContextFactory,
	wire.Bind(
		new(cntxt.IFactory),
		new(*repos.ContextFactory),
	),
	repos.NewForeignsRepository,
	wire.Bind(
		new(foreigns.IRepository),
		new(*repos.ForeignsRepository),
	),
	repos.NewUniquesRepository,
	wire.Bind(
		new(uniques.IRepository),
		new(*repos.UniquesRepository),
	),
	repos.NewUIDRepository,
	wire.Bind(
		new(uids.IRepository),
		new(*repos.UIDRepository),
	),
	repos.NewTasksRepository,
	wire.Bind(
		new(tasks.IRepository),
		new(*repos.TasksRepository),
	),
	repos.NewQuotesRepository,
	wire.Bind(
		new(quotes.IRepository),
		new(*repos.QuotesRepository),
	),

	NewBadTracer,
	wire.Bind(
		new(trace.IRepository),
		new(*badTracer),
	),
)

func NewBadTracer() *badTracer {
	return &badTracer{}
}

type badTracer struct{}

func (_ *badTracer) TraceRequest(
	ctx context.Context,
	method string,
	path string,
	query string,
	statusCode int,
	bodySize int,
	ip string,
	userAgent string,
	startTimestamp time.Time,
	eventTimestamp time.Time,
	fields map[string]string,
) {
}

// Implementation used for graceful starting and stopping of the implementation
// layer
type Implementation struct{}

// NewImplementation constructor for the inmem implementation
func NewImplementation() *Implementation {
	return &Implementation{}
}

// Start runs any routines that are required before the implemtation layer can
// be utilized
func (i *Implementation) Start(ctx context.Context) error {
	return gorr.NewNotImplemented()
}

// Stop runs any routines that are required for the implementation layer to
// gracefully shutdown
func (i *Implementation) Stop(ctx context.Context) error {
	return gorr.NewNotImplemented()
}
