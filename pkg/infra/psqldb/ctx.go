package psqldb

import (
	"fmt"

	"github.com/BetaLixT/tsqlx"
	"github.com/jmoiron/sqlx"
)

func NewDatabaseContext(
	tracer tsqlx.ITracer,
	optn *DatabaseOptions,
) *tsqlx.TracedDB {

	db, err := sqlx.Open("postgres", optn.ConnectionString)
	if err != nil {
		panic(fmt.Errorf("database connection open failure: %w", err))
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("database ping failure"))
	}

	return tsqlx.NewTracedDB(
		db,
		tracer,
		optn.DatabaseServiceName,
	)
}
