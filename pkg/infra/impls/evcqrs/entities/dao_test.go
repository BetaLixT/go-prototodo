package entities

import (
	"context"
	"prototodo/pkg/infra/psqldb"
	"testing"
	"time"
)

type MockTracer struct {
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

}

func TestMigration(t *testing.T) {
	dbctx := psqldb.NewDatabaseContext(
		&MockTracer{},
		&psqldb.DatabaseOptions{
			ConnectionString:    "host=127.0.0.1 port=5433 user=admin password=123456 dbname=todo-test sslmode=disable",
			DatabaseServiceName: "test-db",
		},
	)
	if err := dbctx.Ping(); err != nil {
		println("failed to ping database, skipping test")
		t.SkipNow()
	}
}
