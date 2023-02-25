package repos

import (
	"context"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/cntxt"
	"techunicorn.com/udc-core/prototodo/pkg/domain/common"
	"techunicorn.com/udc-core/prototodo/pkg/domain/domains/tasks"
	"techunicorn.com/udc-core/prototodo/pkg/infra/impls/evcqrs/entities"
	"techunicorn.com/udc-core/prototodo/pkg/infra/psqldb"
	"testing"
	"time"

	"github.com/BetaLixT/tsqlx"
	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
)

type MockTracer struct {
	lgr *zap.Logger
}

func (t *MockTracer) TraceDependency(
	ctx context.Context,
	spanID string,
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
		zap.String("spanId", spanID),
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
	cntxt.IFactory,
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

	return ctxf, lgrf, dbctx, nil
}

func TestCreateDelete(t *testing.T) {
	ctxf, lgrf, dbctx, err := createDependenciesAndMigrate()
	if err != nil {
		println("failed to create dependencies")
		t.SkipNow()
	}

	base := NewBaseDataRepository(dbctx)
	r := NewTasksRepository(
		base,
		lgrf,
	)

	ctx := ctxf.Create(
		context.Background(),
		time.Minute*2,
	)

	lgr := lgrf.Create(ctx)
	sf, err := snowflake.NewNode(1)
	if err != nil {
		lgr.Error("failed to create snowflake", zap.Error(err))
	}

	id := sf.Generate().String()
	ev, err := r.Create(
		ctx,
		id,
		nil,
		tasks.TaskData{
			Title:       Pointerify("title"),
			Description: Pointerify("description"),
		},
	)
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

func TestCreateUpdate(t *testing.T) {
	ctxf, lgrf, dbctx, err := createDependenciesAndMigrate()
	if err != nil {
		println("failed to create dependencies")
		t.SkipNow()
	}

	base := NewBaseDataRepository(dbctx)
	r := NewTasksRepository(
		base,
		lgrf,
	)

	ctx := ctxf.Create(
		context.Background(),
		time.Minute*2,
	)

	lgr := lgrf.Create(ctx)
	sf, err := snowflake.NewNode(1)
	if err != nil {
		lgr.Error("failed to create snowflake", zap.Error(err))
	}

	id := sf.Generate().String()
	_, err = r.Create(
		ctx,
		id,
		nil,
		tasks.TaskData{
			Title:       Pointerify("original title"),
			Description: Pointerify("description"),
		},
	)
	if err != nil {
		println()
		lgr.Error("failed to create record", zap.Error(err))
		t.FailNow()
	}
	err = ctx.CommitTransaction()
	if err != nil {
		lgr.Error("failed commit transaction record", zap.Error(err))
		t.FailNow()
	}

	ctx2 := ctxf.Create(
		context.Background(),
		time.Minute*2,
	)

	ev, err := r.Update(
		ctx2,
		id,
		nil,
		1,
		tasks.TaskData{
			Title: Pointerify("updated title"),
		},
	)
	if err != nil {
		lgr.Error("failed to update record", zap.Error(err))
		t.FailNow()
	}
	err = ctx2.CommitTransaction()
	if err != nil {
		lgr.Error("failed commit transaction record", zap.Error(err))
		t.FailNow()
	}

	if *ev.Data.Title != "updated title" {
		lgr.Error("invalid title from event", zap.Stringp("title", ev.Data.Title))
		t.FailNow()
	}

	ctx3 := ctxf.Create(
		context.Background(),
		time.Minute*2,
	)
	v, err := r.Get(ctx3, id)
	if err != nil {
		lgr.Error("failed to get task")
		t.FailNow()
	}
	if v.Title != "updated title" {
		lgr.Error("invalid title on read model", zap.String("title", v.Title))
	}
	if v.Description != "description" {
		lgr.Error("invalid description on read model", zap.String("desc", v.Description))
	}
	if v.Version != 1 {
		lgr.Error("invalid version on read model", zap.Uint64("version", v.Version))
	}
}

func TestFullUpdate(t *testing.T) {
	ctxf, lgrf, dbctx, err := createDependenciesAndMigrate()
	if err != nil {
		println("failed to create dependencies")
		t.SkipNow()
	}

	base := NewBaseDataRepository(dbctx)
	r := NewTasksRepository(
		base,
		lgrf,
	)

	ctx := ctxf.Create(
		context.Background(),
		time.Minute*2,
	)

	lgr := lgrf.Create(ctx)
	sf, err := snowflake.NewNode(1)
	if err != nil {
		lgr.Error("failed to create snowflake", zap.Error(err))
	}

	id := sf.Generate().String()
	_, err = r.Create(
		ctx,
		id,
		nil,
		tasks.TaskData{
			Title:       Pointerify("original title"),
			Description: Pointerify("description"),
		},
	)
	if err != nil {
		println()
		lgr.Error("failed to create record", zap.Error(err))
		t.FailNow()
	}
	err = ctx.CommitTransaction()
	if err != nil {
		lgr.Error("failed commit transaction record", zap.Error(err))
		t.FailNow()
	}

	ctx2 := ctxf.Create(
		context.Background(),
		time.Minute*2,
	)

	_, err = r.Update(
		ctx2,
		id,
		nil,
		1,
		tasks.TaskData{
			Title:       Pointerify("updated title"),
			Description: Pointerify("updated description"),
			Status:      Pointerify("completed"),
			RandomMap:   map[string]string{"wow": "wowie"},
			Metadata:    map[string]interface{}{"cries": 34},
		},
	)
	if err != nil {
		lgr.Error("failed to update record", zap.Error(err))
		t.FailNow()
	}
	err = ctx2.CommitTransaction()
	if err != nil {
		lgr.Error("failed commit transaction record", zap.Error(err))
		t.FailNow()
	}

	ctx3 := ctxf.Create(
		context.Background(),
		time.Minute*2,
	)
	v, err := r.Get(ctx3, id)
	if err != nil {
		lgr.Error("failed to get task")
		t.FailNow()
	}
	if v.Title != "updated title" {
		lgr.Error("invalid title from event", zap.String("title", v.Title))
		t.FailNow()
	}
	if v.Description != "updated description" {
		lgr.Error("invalid description from event", zap.String("description", v.Description))
		t.FailNow()
	}
	if v.Status != "completed" {
		lgr.Error("invalid status from event", zap.String("status", v.Status))
		t.FailNow()
	}
	if v.RandomMap["wow"] != "wowie" {
		lgr.Error("invalid map from event", zap.Any("randomMap", v.RandomMap))
		t.FailNow()
	}
	mval := v.Metadata["cries"]
	if mval.(float64) != 34 {
		lgr.Error("invalid meta from event", zap.Any("status", v.Metadata))
		t.FailNow()
	}
}

func Pointerify[x any](val x) *x { return &val }
