// Package inmem Event source CQRS implementation of the domain layer
package inmem

import (
	"context"
	"prototodo/pkg/domain/base/acl"
	"prototodo/pkg/domain/base/foreigns"
	"prototodo/pkg/domain/base/impl"
	"prototodo/pkg/domain/base/logger"
	"prototodo/pkg/domain/base/tracectx"
	"prototodo/pkg/domain/base/uids"
	"prototodo/pkg/domain/base/uniques"
	"prototodo/pkg/domain/domains/quotes"
	"prototodo/pkg/domain/domains/tasks"
	"prototodo/pkg/infra/config"
	"prototodo/pkg/infra/impls/inmem/repos"
	"prototodo/pkg/infra/lgr"
	"prototodo/pkg/infra/sf"

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
	repos.NewForeignsRepository,
	wire.Bind(
		new(foreigns.IRepository),
		new(*repos.ForeignsRepository),
	),
	repos.NewUniquesRepository,
	wire.Bind(
		new(uniques.IRepository),
		new(repos.UniquesRepository),
	),
	repos.NewTraceRepository,
	wire.Bind(
		new(tracectx.IRepository),
		new(repos.TraceRepository),
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
)

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
