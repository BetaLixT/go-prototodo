package repos

import (
	"context"
	"prototodo/pkg/domain/base/acl"
	"testing"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const (
	TestRedisConnString = ""
)

func TestCreateACLEntry(t *testing.T) {
	ctxf, lgrf, dbctx, err := createDependenciesAndMigrate()
	if err != nil {
		println("failed to create dependencies")
		t.SkipNow()
	}
	lgr := lgrf.Create(context.Background())

	rdb := redis.NewClient(
		&redis.Options{
			Addr: "127.0.0.1:6379",
			DB:   0,
		},
	)
	err = rdb.Ping(context.Background()).Err()
	if err != nil {
		println("failed creating redis connection")
		t.SkipNow()
	}

	base := NewBaseDataRepository(dbctx)
	r := NewACLRepository(
		base,
		rdb,
		lgrf,
	)

	sf, err := snowflake.NewNode(1)
	if err != nil {
		lgr.Error("failed to create snowflake", zap.Error(err))
	}

	id := sf.Generate().String()

	ctx1 := ctxf.Create(
		context.Background(),
		time.Minute*2,
	)
	err = r.CreateACLEntry(
		ctx1,
		id,
		"123",
		"tester",
		"123",
		acl.Read,
	)
	if err != nil {
		lgr.Error("acl creation failed", zap.Error(err))
		t.FailNow()
	}
	ctx1.RollbackTransaction()

	ctxr := ctxf.Create(
		context.Background(),
		time.Minute*2,
	)

	err = r.CanRead(
		ctxr,
		id,
		[]string{"123"},
		"tester",
		"xyz",
	)
	if err == nil {
		lgr.Error("expected an error, but no errors")
		t.FailNow()
	}

	ctx2 := ctxf.Create(
		context.Background(),
		time.Minute*2,
	)
	err = r.CreateACLEntry(
		ctx2,
		id,
		"123",
		"tester",
		"xyz",
		acl.Read,
	)
	if err != nil {
		lgr.Error("acl creation failed", zap.Error(err))
		t.FailNow()
	}
	ctx2.CommitTransaction()

	err = r.CanRead(
		ctxr,
		id,
		[]string{"123"},
		"tester",
		"xyz",
	)
	if err != nil {
		lgr.Error("expected can read but failed", zap.Error(err))
		t.FailNow()
	}

	// running twice to test caching
	err = r.CanRead(
		ctxr,
		id,
		[]string{"123"},
		"tester",
		"xyz",
	)
	if err != nil {
		lgr.Error("expected can read but failed second time around", zap.Error(err))
		t.FailNow()
	}

	err = r.CanWrite(
		ctxr,
		id,
		[]string{"123"},
		"tester",
		"xyz",
	)
	if err == nil {
		lgr.Error("expected failure for can write cries")
		t.FailNow()
	}
}
