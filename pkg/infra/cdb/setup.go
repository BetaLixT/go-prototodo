package cdb

import (
	"context"
	"fmt"

	"github.com/gocql/gocql"
	"go.uber.org/zap"
	"techunicorn.com/udc/lette/pkg/infra/logger"
)

type CassandraSetup struct {
	sess *gocql.Session
	lgrf *logger.LoggerFactory
}

func NewCassandraSetup(
	sess *gocql.Session,
	lgrf *logger.LoggerFactory,
) *CassandraSetup {
	return &CassandraSetup{
		sess: sess,
		lgrf: lgrf,
	}
}

func (cs *CassandraSetup) Initialize(
	ctx context.Context,
) error {
	lgr := cs.lgrf.NewLogger(ctx)
	err := cs.sess.Query(Schema).Exec()
	if err != nil {
		lgr.Error("error while initializing cassandra", zap.Error(err))
		return fmt.Errorf("error while initializing cassandra: %w", err)
	}
	return nil
}

const (
	Schema = `
  CREATE TABLE IF NOT EXISTS livestream_user (
    session_id text,
    participant_id text,
    id uuid,
    PRIMARY KEY ((session_id, participant_id))
  );
  `
)
