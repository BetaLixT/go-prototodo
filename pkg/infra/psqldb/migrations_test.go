package psqldb

import (
	"context"
	"testing"
	"time"

	"github.com/BetaLixT/tsqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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

func TestRunMigrations(t *testing.T) {
	ctx := context.TODO()

	lgr, _ := zap.NewProduction()
	viper.SetConfigName("config")
	viper.KeyDelimiter("__")
	viper.AddConfigPath("../../../cfg")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			lgr.Warn("No config file found")
		} else {
			lgr.Error("Failed loading config")
			t.FailNow()
		}
	}

	dbctx := NewDatabaseContext(
		&MockTracer{},
		&DatabaseOptions{
			ConnectionString: viper.GetString("DatabaseOptions.ConnectionString"),
		})

	err := RunMigrations(
		ctx,
		lgr,
		tsqlx.NewTracedDB(
			dbctx.DB,
			&MockTracer{},
			"main-database",
		),
		GetMigrationScripts(),
	)
	if err != nil {
		lgr.Error("Failed migrations", zap.Error(err))
		t.FailNow()
	}
}
