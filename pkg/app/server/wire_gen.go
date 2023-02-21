// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"prototodo/pkg/app/server/handlers"
	"prototodo/pkg/domain/domains/quotes"
	"prototodo/pkg/domain/domains/tasks"
	"prototodo/pkg/infra/config"
	"prototodo/pkg/infra/impls/evcqrs"
	"prototodo/pkg/infra/impls/evcqrs/repos"
	"prototodo/pkg/infra/impls/inmem"
	repos2 "prototodo/pkg/infra/impls/inmem/repos"
	"prototodo/pkg/infra/lgr"
	"prototodo/pkg/infra/psqldb"
	"prototodo/pkg/infra/rdb"
	"prototodo/pkg/infra/sf"
	"prototodo/pkg/infra/trace"
	"prototodo/pkg/infra/trace/appinsights"
	"prototodo/pkg/infra/trace/jaeger"
	"prototodo/pkg/infra/trace/promex"
)

// Injectors from wire.go:

// InitializeEvent creates an Event. It will error if the Event is staffed with
// a grumpy greeter.
func initializeAppCQRS() (*app, error) {
	loggerFactory, err := lgr.NewLoggerFactory()
	if err != nil {
		return nil, err
	}
	initializer := config.NewInitializer(loggerFactory)
	exporterOptions := config.NewAppInsightsExporterOptions(initializer)
	traceExporter, err := appinsights.NewTraceExporter(exporterOptions)
	if err != nil {
		return nil, err
	}
	jaegerExporterOptions := config.NewJaegerExporterOptions(initializer)
	jaegerTraceExporter, err := jaeger.NewJaegerTraceExporter(jaegerExporterOptions)
	if err != nil {
		return nil, err
	}
	promexTraceExporter, err := promex.NewTraceExporter()
	if err != nil {
		return nil, err
	}
	exporterList := evcqrs.NewTraceExporterList(traceExporter, jaegerTraceExporter, promexTraceExporter, loggerFactory)
	options, err := config.NewTraceOptions(initializer)
	if err != nil {
		return nil, err
	}
	tracer, err := trace.NewTracer(exporterList, options, loggerFactory)
	if err != nil {
		return nil, err
	}
	databaseOptions := config.NewPSQLDBOptions(initializer)
	tracedDB, err := psqldb.NewDatabaseContext(tracer, databaseOptions)
	if err != nil {
		return nil, err
	}
	baseDataRepository := repos.NewBaseDataRepository(tracedDB)
	tasksRepository := repos.NewTasksRepository(baseDataRepository, loggerFactory)
	rdbOptions := config.NewRedisOptions(initializer)
	client, err := rdb.NewRedisContext(rdbOptions, tracer)
	if err != nil {
		return nil, err
	}
	aclRepository := repos.NewACLRepository(baseDataRepository, client, loggerFactory)
	sfOptions := config.NewSnowflakeOptions(initializer)
	node, err := sf.NewSnowflake(sfOptions)
	if err != nil {
		return nil, err
	}
	uidRepository := repos.NewUIDRepository(node)
	service := tasks.NewService(tasksRepository, loggerFactory, aclRepository, uidRepository)
	tasksHandler := handlers.NewTasksHandler(loggerFactory, service)
	quotesRepository := repos.NewQuotesRepository(tracedDB, loggerFactory)
	quotesService := quotes.NewService(quotesRepository, loggerFactory, uidRepository)
	quotesHandler := handlers.NewQuotesHandler(loggerFactory, quotesService)
	implementation := evcqrs.NewImplementation(tracedDB, loggerFactory)
	contextFactory := repos.NewContextFactory(loggerFactory)
	serverApp := newApp(tasksHandler, quotesHandler, tasksHandler, quotesHandler, implementation, loggerFactory, contextFactory)
	return serverApp, nil
}

func initializeAppInMem() (*app, error) {
	loggerFactory, err := lgr.NewLoggerFactory()
	if err != nil {
		return nil, err
	}
	tasksRepository := repos2.NewTasksRepository()
	aclRepository := repos2.NewACLRepository()
	initializer := config.NewInitializer(loggerFactory)
	options := config.NewSnowflakeOptions(initializer)
	node, err := sf.NewSnowflake(options)
	if err != nil {
		return nil, err
	}
	uidRepository := repos2.NewUIDRepository(node)
	service := tasks.NewService(tasksRepository, loggerFactory, aclRepository, uidRepository)
	tasksHandler := handlers.NewTasksHandler(loggerFactory, service)
	quotesRepository := repos2.NewQuotesRepository()
	quotesService := quotes.NewService(quotesRepository, loggerFactory, uidRepository)
	quotesHandler := handlers.NewQuotesHandler(loggerFactory, quotesService)
	implementation := inmem.NewImplementation()
	contextFactory := repos2.NewContextFactory()
	serverApp := newApp(tasksHandler, quotesHandler, tasksHandler, quotesHandler, implementation, loggerFactory, contextFactory)
	return serverApp, nil
}
