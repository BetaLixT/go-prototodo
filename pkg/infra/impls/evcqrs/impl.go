package evcqrs

import (
	"context"
	"prototodo/pkg/domain/base/acl"
	"prototodo/pkg/domain/base/foreigns"
	"prototodo/pkg/domain/base/logger"
	"prototodo/pkg/domain/base/trace"
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

	"github.com/BetaLixT/tsqlx"
	"github.com/google/wire"
	"go.uber.org/zap"
)

var DependencySet = wire.NewSet(
	// Infra
	lgr.NewLoggerFactory,
	wire.Bind(
		new(logger.IFactory),
		new(*lgr.LoggerFactory),
	),
	psqldb.NewDatabaseContext,
	config.NewPSQLDBOptions,
	rdb.NewRedisContext,
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
		new(trace.IRepository),
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

type Implementation struct {
	dbctx *tsqlx.TracedDB
	lgrf  *lgr.LoggerFactory
}

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

func (i *Implementation) Stop(ctx context.Context) error {
	i.lgrf.Close()
	return nil
}
