// Package server contains server logic to handle incoming requests and command
// query handlers
package server

import (
	"context"
	"os"
	"os/signal"
	"prototodo/pkg/app/server/contracts"
	"prototodo/pkg/app/server/handlers"
	"prototodo/pkg/domain"
	"prototodo/pkg/domain/base/cntxt"
	"prototodo/pkg/domain/base/impl"
	"prototodo/pkg/domain/base/logger"
	"prototodo/pkg/infra/impls/evcqrs"
	"prototodo/pkg/infra/impls/inmem"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Start boots up the server
func Start(impl string) {
	var a *app
	var err error
	switch impl {
	case "inmem":
		a, err = initializeAppInMem()
		if err != nil {
			panic(err)
		}
	default:
		a, err = initializeAppCQRS()
		if err != nil {
			panic(err)
		}
	}

	a.start(context.Background())
}

// cqrsDependencySet dependency set with in memory CQRS implementation
var cqrsDependencySet = wire.NewSet(
	evcqrs.DependencySet,
	domain.DependencySet,
	dependencySet,
)

// inMemDependencySet dependency set with in memory implementation
var inMemDependencySet = wire.NewSet(
	inmem.DependencySet,
	domain.DependencySet,
	dependencySet,
)

var dependencySet = wire.NewSet(
	newApp,
	handlers.NewQuotesHandler,
	wire.Bind(
		new(contracts.QuotesHTTPServer),
		new(*handlers.QuotesHandler),
	),
	wire.Bind(
		new(contracts.QuotesServer),
		new(*handlers.QuotesHandler),
	),

	handlers.NewTasksHandler,
	wire.Bind(
		new(contracts.TasksHTTPServer),
		new(*handlers.TasksHandler),
	),
	wire.Bind(
		new(contracts.TasksServer),
		new(*handlers.TasksHandler),
	),
)

// =============================================================================
// Application
// =============================================================================

type closer func()

type app struct {
	// http handler interfaces
	tasksHTTPHandler  contracts.TasksHTTPServer
	quotesHTTPHandler contracts.QuotesHTTPServer

	// grpc handler interfaces
	tasksGRPCHandler  contracts.TasksServer
	quotesGRPCHandler contracts.QuotesServer

	impl impl.IImplementation
	lgr  *zap.Logger
	ctxf cntxt.IFactory

	// server closers
	closers   []closer
	closeLock sync.Mutex
}

func newApp(
	tasksHTTPHandler contracts.TasksHTTPServer,
	quotesHTTPHandler contracts.QuotesHTTPServer,
	tasksGRPCHandler contracts.TasksServer,
	quotesGRPCHandler contracts.QuotesServer,
	impl impl.IImplementation,
	lgrf logger.IFactory,
	ctxf cntxt.IFactory,
) *app {
	return &app{
		// http handler interfaces
		tasksHTTPHandler:  tasksHTTPHandler,
		quotesHTTPHandler: quotesHTTPHandler,

		// grpc handler interfaces
		tasksGRPCHandler:  tasksGRPCHandler,
		quotesGRPCHandler: quotesGRPCHandler,

		impl: impl,
		lgr:  lgrf.Create(context.Background()),
		ctxf: ctxf,
	}
}

func (a *app) registerGRPCHandlers(s *grpc.Server) {
	contracts.RegisterTasksServer(s, a.tasksGRPCHandler)
	contracts.RegisterQuotesServer(s, a.quotesGRPCHandler)
}

func (a *app) registerHTTPHandlers(g *gin.RouterGroup) {
	contracts.RegisterTasksHTTPServer(g, a.tasksHTTPHandler)
	contracts.RegisterQuotesHTTPServer(g, a.quotesHTTPHandler)
}

func (a *app) start(ctx context.Context) {
	err := a.impl.Start(ctx)
	if err != nil {
		a.lgr.Error("failed to start implementation", zap.Error(err))
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		a.startGRPC(os.Getenv("PORT_GRPC"))
		a.lgr.Info("grpc server closing...")
	}()

	go func() {
		defer wg.Done()
		a.startHTTP(os.Getenv("PORT_HTTP"))
		a.lgr.Info("http server closing...")
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling

	a.lgr.Info("Server shutting down...")
	a.closeServers()

	wg.Wait()
	a.impl.Stop(ctx)
	a.lgr.Info("Server exiting")
}

func (a *app) registerCloser(c closer) {
	a.closeLock.Lock()
	a.closers = append(a.closers, c)
	a.closeLock.Unlock()
}

func (a *app) closeServers() {
	a.closeLock.Lock()
	for idx := range a.closers {
		a.closers[idx]()
	}
	a.closeLock.Unlock()
}
