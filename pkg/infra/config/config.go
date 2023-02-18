package config

import (
	"context"
	"os"
	"prototodo/pkg/domain/base/logger"
	"prototodo/pkg/infra/cdb"
	"prototodo/pkg/infra/psqldb"
	"prototodo/pkg/infra/rdb"
	"prototodo/pkg/infra/sf"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type ConfigInitializer struct {
	lgrf logger.IFactory
}

// - this is very very scuffed lmao
func NewConfigInitializer(
	lgrf logger.IFactory,
) *ConfigInitializer {
	c := &ConfigInitializer{
		lgrf: lgrf,
	}
	c.loadConfig()
	return c
}

func (c *ConfigInitializer) loadConfig() {
	c.LoadConfigCustom("./cfg/.env")
}

func (c *ConfigInitializer) LoadConfigCustom(loc string) {
	err := godotenv.Load(loc)
	if err != nil {
		lgr := c.lgrf.Create(context.Background())
		lgr.Warn(
			"failed to load .env file, probably missing",
			zap.String("location", loc),
		)
	}
}

func NewCassandraOptions(c *ConfigInitializer) *cdb.Options {
	ips := os.Getenv("CassandraClusterIPs")
	opt := &cdb.Options{
		ClusterIPs: strings.Split(ips, ","),
		Username:   os.Getenv("CassandraUsername"),
		Password:   os.Getenv("CassandraPassword"),
		Keyspace:   os.Getenv("CassandraKeyspace"),
	}
	if len(opt.ClusterIPs) == 0 {
		panic("cassandra url not provided")
	}
	if opt.Keyspace == "" {
		panic("cassandra keyspace not provided")
	}
	lgr := c.lgrf.Create(context.Background())
	if opt.Username == "" {
		lgr.Warn("cassandra options, username missing")
	}
	if opt.Password == "" {
		lgr.Warn("cassandra options, password missing")
	}
	return opt
}

func NewRedisOptions(c *ConfigInitializer) *rdb.Options {
	address := os.Getenv("RedisAddress")
	if address == "" {
		panic("missing redis address config")
	}
	password := os.Getenv("RedisPassword")
	if password == "" {
		lgr := c.lgrf.Create(context.Background())
		lgr.Warn("redis password missing")
	}
	tls := os.Getenv("RedisTls") == "true"
	databaseNumber := os.Getenv("RedisDatabase")
	db, err := strconv.Atoi(databaseNumber)
	if err != nil {
		db = 0
		lgr := c.lgrf.Create(context.Background())
		lgr.Warn("no database number was provided for redis, using default")
	}

	return &rdb.Options{
		Address:     address,
		Password:    password,
		ServiceName: address,
		TLS:         tls,
		Database:    db,
	}
}

func NewSnowflakeOptions(_ *ConfigInitializer) *sf.Options {
	nn := os.Getenv("SnowflakeNodeNumber")
	if nn == "" {
		panic("missing snowflake node number")
	}
	nni, err := strconv.ParseInt(
		nn,
		10,
		64,
	)
	if err != nil {
		panic("failed to parse node number")
	}

	return &sf.Options{
		NodeNumber: nni,
	}
}

func NewPSQLDBOptions(_ *ConfigInitializer) *psqldb.DatabaseOptions {
	cons := os.Getenv("DatabaseConnectionString")
	if cons == "" {
		panic("missing database connection string config")
	}

	split := strings.Split(cons, " ")
	name := "psql-database"
	for idx := range split {
		if strings.HasPrefix(split[idx], "host=") {
			name = strings.TrimPrefix(split[idx], "host=")
		}
	}
	return &psqldb.DatabaseOptions{
		ConnectionString:    cons,
		DatabaseServiceName: name,
	}
}
