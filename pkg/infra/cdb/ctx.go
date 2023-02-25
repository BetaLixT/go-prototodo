package cdb

import (
	"context"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/logger"
	"strconv"

	trace "github.com/BetaLixT/appInsightsTrace"
	"github.com/gocql/gocql"
	"go.uber.org/zap"
)

func NewCassandraSession(
	optn *Options,
	lgrf logger.IFactory,
	isn *trace.AppInsightsCore,
) (*gocql.Session, error) {
	lgr := lgrf.Create(context.Background())
	lgr.Info("connecting to cassandra...", zap.Strings("urls", optn.ClusterIPs))
	auth := gocql.PasswordAuthenticator{
		Username: optn.Username,
		Password: optn.Password,
	}

	// - creating main session
	cls := gocql.NewCluster(optn.ClusterIPs...)
	cls.Authenticator = auth
	cls.Keyspace = optn.Keyspace
	cls.QueryObserver = &TraceObserver{
		ins: isn,
	}
	sess, err := cls.CreateSession()
	if err != nil {
		return nil, err
	}

	return sess, nil
}

type TraceObserver struct {
	ins *trace.AppInsightsCore
}

func (o *TraceObserver) ObserveQuery(ctx context.Context, qry gocql.ObservedQuery) {
	o.ins.TraceDependency(
		ctx,
		"",
		"cassandra",
		qry.Host.ClusterName(),
		qry.Statement,
		qry.Err == nil,
		qry.Start,
		qry.End,
		map[string]string{
			"keyspace":    qry.Keyspace,
			"serverLtncy": strconv.FormatInt(qry.Metrics.TotalLatency, 10),
		},
	)
}

const (
	CreateKeyspace = `
  CREATE KEYSPACE IF NOT EXISTS %s WITH REPLICATION = {
  	'class' : 'SimpleStrategy',
  	'replication_factor' : '1'
  };
  `
)
