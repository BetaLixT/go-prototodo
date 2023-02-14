package repos

import (
	"context"
	"prototodo/pkg/domain/base/cntxt"
	"prototodo/pkg/domain/common"
	"prototodo/pkg/domain/domains/tasks"
	"prototodo/pkg/infra/impls/evcqrs/entities"
	"prototodo/pkg/infra/psqldb"
	"testing"
	"time"

	"github.com/BetaLixT/tsqlx"
	"go.uber.org/zap"
)

type MockTracer struct {
	lgr *zap.Logger
}

func (t *MockTracer) TraceDependency(
	ctx context.Context,
	spanId string,
	dependencyType string,
	serviceName string,
	commandName string,
	success bool,
	startTimestamp time.Time,
	eventTimestamp time.Time,
	fields map[string]string,
) {
	t.lgr.Info(
		"dependency",
		zap.String("spanId", spanId),
		zap.String("dependencyType", dependencyType),
		zap.String("serviceName", serviceName),
		zap.String("commandName", commandName),
		zap.Bool("success", success),
		zap.Any("fields", fields),
	)
}

type LoggerFactory struct {
	lgr *zap.Logger
}

func (f *LoggerFactory) Create(_ context.Context) *zap.Logger {
	return f.lgr
}

const (
	TestDatabaseConnString  = "host=127.0.0.1 port=5432 user=admin password=123456 dbname=todo_test sslmode=disable"
	TestDatabaseServiceName = "test-db"
)

func createDependenciesAndMigrate() (
	cntxt.IContext,
	*LoggerFactory,
	*tsqlx.TracedDB,
	error,
) {
	lgr, _ := zap.NewDevelopment()
	dbctx, err := psqldb.NewDatabaseContext(
		&MockTracer{lgr: lgr},
		&psqldb.DatabaseOptions{
			ConnectionString:    TestDatabaseConnString,
			DatabaseServiceName: TestDatabaseServiceName,
		},
	)
	if err != nil {
		println("failed to connect to database")
		return nil, nil, nil, err
	}

	lgrf := &LoggerFactory{lgr: lgr}

	ctxf := NewContextFactory(
		lgrf,
		NewTraceRepository(lgrf),
	)

	ctx := ctxf.Create(
		context.Background(),
		time.Minute*2,
	)
	err = psqldb.RunMigrations(
		ctx,
		lgr,
		dbctx,
		entities.GetMigrationScripts(),
	)
	if err != nil {
		println("failed to run migration")
		return nil, nil, nil, err
	}

	return ctx, lgrf, dbctx, nil
}

func TestCreateDelete(t *testing.T) {
	ctx, lgrf, dbctx, err := createDependenciesAndMigrate()
	if err != nil {
		println("failed to create dependencies")
		t.SkipNow()
	}

	r := NewTasksRepository(
		dbctx,
		lgrf,
	)

	id := "1"
	ev, err := r.Create(
		ctx,
		id,
		nil,
		tasks.TaskData{
			Title:       Pointerify("title"),
			Description: Pointerify("description"),
		},
	)
	lgr := lgrf.Create(ctx)
	if err != nil {
		println()
		lgr.Error("failed to create record", zap.Error(err))
		t.FailNow()
	}
	if ev.StreamId != id {
		println("invalid id")
		t.FailNow()
	}
	if ev.Stream != common.TaskStreamName {
		println("invalid stream")
		t.FailNow()
	}
	if *ev.Data.Title != "title" {
		println("invalid title")
		t.FailNow()
	}
	_, err = r.Delete(
		ctx,
		id,
		nil,
		1,
	)
	if err != nil {
		lgr.Error("failed to delete record", zap.Error(err))
		t.FailNow()
	}
	err = ctx.CommitTransaction()
	if err != nil {
		lgr.Error("failed to commit transaction", zap.Error(err))
		t.FailNow()
	}
}

func Pointerify[x comparable](val x) *x { return &val }
