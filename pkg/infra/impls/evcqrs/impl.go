// Package evcqrs Event source CQRS implementation of the domain layer
package evcqrs

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
	"prototodo/pkg/infra/impls/evcqrs/entities"
	"prototodo/pkg/infra/impls/evcqrs/repos"
	"prototodo/pkg/infra/lgr"
	"prototodo/pkg/infra/psqldb"
	"prototodo/pkg/infra/rdb"
	"prototodo/pkg/infra/sf"
	"prototodo/pkg/infra/trace"
	"prototodo/pkg/infra/trace/appinsights"
	"prototodo/pkg/infra/trace/jaeger"
	"prototodo/pkg/infra/tracelib"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/BetaLixT/gotred/v8"
	"github.com/BetaLixT/tsqlx"
	"github.com/google/wire"
	"go.uber.org/zap"
)

// DependencySet dependencies provided by the implementation
var DependencySet = wire.NewSet(
	NewImplementation,
	wire.Bind(
		new(impl.IImplementation),
		new(*Implementation),
	),
	// Trace
	NewTraceExporterList,
	config.NewTraceOptions,
	trace.NewTracer,

	// Infra
	lgr.NewLoggerFactory,
	wire.Bind(
		new(logger.IFactory),
		new(*lgr.LoggerFactory),
	),
	psqldb.NewDatabaseContext,
	wire.Bind(
		new(tsqlx.ITracer),
		new(*tracelib.Tracer),
	),
	config.NewPSQLDBOptions,
	rdb.NewRedisContext,
	wire.Bind(
		new(gotred.ITracer),
		new(*tracelib.Tracer),
	),
	config.NewRedisOptions,
	sf.NewSnowflake,
	config.NewSnowflakeOptions,

	// Repos
	repos.NewBaseDataRepository,
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

// NewTraceExporterList provides a list of exporters for tracing
func NewTraceExporterList(
	insexp appinsights.TraceExporter,
	jgrexp jaeger.TraceExporter,
	lgrf logger.IFactory,
) *trace.ExporterList {
	lgr := lgrf.Create(context.Background())
	exp := []sdktrace.SpanExporter{}

	if insexp != nil {
		exp = append(exp, insexp)
	} else {
		lgr.Warn("insights exporter not found")
	}
	if jgrexp != nil {
		exp = append(exp, jgrexp)
	} else {
		lgr.Warn("jeager exporter not found")
	}
	if len(exp) == 0 {
		panic("no tracing exporters found (float you <3)")
	}
	return &trace.ExporterList{
		Exporters: exp,
	}
}

// Implementation used for graceful starting and stopping of the implementation
// layer
type Implementation struct {
	dbctx *tsqlx.TracedDB
	lgrf  *lgr.LoggerFactory
}

// NewImplementation constructor for the evcqrs implementation
func NewImplementation(
	dbctx *tsqlx.TracedDB,
	lgrf *lgr.LoggerFactory,
) *Implementation {
	return &Implementation{
		dbctx: dbctx,
		lgrf:  lgrf,
	}
}

// Start runs any routines that are required before the implemtation layer can
// be utilized
func (i *Implementation) Start(ctx context.Context) error {
	lgri := i.lgrf.Create(ctx)
	err := psqldb.RunMigrations(
		ctx,
		lgri,
		i.dbctx,
		entities.GetMigrationScripts(),
	)
	if err != nil {
		lgri.Error("failed to run migration", zap.Error(err))
		return err
	}
	return nil
}

// Stop runs any routines that are required for the implementation layer to
// gracefully shutdown
func (i *Implementation) Stop(ctx context.Context) error {
	i.lgrf.Close()
	return nil
}
