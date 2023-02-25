// Package evcqrs Event source CQRS implementation of the domain layer
package evcqrs

import (
	"context"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/acl"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/cntxt"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/foreigns"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/impl"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/logger"
	domtrace "techunicorn.com/udc-core/prototodo/pkg/domain/base/trace"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/uids"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/uniques"
	"techunicorn.com/udc-core/prototodo/pkg/domain/domains/quotes"
	"techunicorn.com/udc-core/prototodo/pkg/domain/domains/tasks"
	"techunicorn.com/udc-core/prototodo/pkg/infra/config"
	"techunicorn.com/udc-core/prototodo/pkg/infra/impls/evcqrs/entities"
	"techunicorn.com/udc-core/prototodo/pkg/infra/impls/evcqrs/repos"
	"techunicorn.com/udc-core/prototodo/pkg/infra/lgr"
	"techunicorn.com/udc-core/prototodo/pkg/infra/psqldb"
	"techunicorn.com/udc-core/prototodo/pkg/infra/rdb"
	"techunicorn.com/udc-core/prototodo/pkg/infra/sf"
	"techunicorn.com/udc-core/prototodo/pkg/infra/trace"
	"techunicorn.com/udc-core/prototodo/pkg/infra/trace/appinsights"
	"techunicorn.com/udc-core/prototodo/pkg/infra/trace/jaeger"
	"techunicorn.com/udc-core/prototodo/pkg/infra/trace/promex"
	"techunicorn.com/udc-core/prototodo/pkg/infra/tracelib"

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
	jaeger.NewJaegerTraceExporter,
	config.NewJaegerExporterOptions,
	appinsights.NewTraceExporter,
	config.NewAppInsightsExporterOptions,
	promex.NewTraceExporter,

	// Infra
	config.NewInitializer,
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

	wire.Bind(
		new(domtrace.IRepository),
		new(*tracelib.Tracer),
	),
)

// NewTraceExporterList provides a list of exporters for tracing
func NewTraceExporterList(
	insexp appinsights.TraceExporter,
	jgrexp jaeger.TraceExporter,
	prmex promex.TraceExporter,
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
	exp = append(exp, prmex)
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
