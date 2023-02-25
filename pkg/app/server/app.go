// Package server contains server logic to handle incoming requests and command
// query handlers
package server

import (
	"context"
	"os"
	"os/signal"
	"techunicorn.com/udc-core/prototodo/pkg/app/server/contracts"
	"techunicorn.com/udc-core/prototodo/pkg/app/server/handlers"
	"techunicorn.com/udc-core/prototodo/pkg/domain"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/cntxt"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/impl"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/logger"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/trace"
	"techunicorn.com/udc-core/prototodo/pkg/infra/impls/evcqrs"
	"techunicorn.com/udc-core/prototodo/pkg/infra/impls/inmem"
	"sync"
	"syscall"
	"time"

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
	lgrf logger.IFactory
	lgr  *zap.Logger
	ctxf cntxt.IFactory
	trc  trace.IRepository

	// server closers
	closers   []closer
	closeLock sync.Mutex
	closed    bool
	quit      chan os.Signal
}

func newApp(
	tasksHTTPHandler contracts.TasksHTTPServer,
	quotesHTTPHandler contracts.QuotesHTTPServer,
	tasksGRPCHandler contracts.TasksServer,
	quotesGRPCHandler contracts.QuotesServer,
	impl impl.IImplementation,
	lgrf logger.IFactory,
	ctxf cntxt.IFactory,
	trc trace.IRepository,
) *app {
	return &app{
		// http handler interfaces
		tasksHTTPHandler:  tasksHTTPHandler,
		quotesHTTPHandler: quotesHTTPHandler,

		// grpc handler interfaces
		tasksGRPCHandler:  tasksGRPCHandler,
		quotesGRPCHandler: quotesGRPCHandler,

		impl: impl,
		lgrf: lgrf,
		lgr:  lgrf.Create(context.Background()),
		ctxf: ctxf,
		trc:  trc,

		quit: make(chan os.Signal, 1),
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
		a.closeServers()
	}()

	go func() {
		defer wg.Done()
		a.startHTTP(os.Getenv("PORT_HTTP"))
		a.lgr.Info("http server closing...")
		a.closeServers()
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	a.quit = make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(a.quit, syscall.SIGINT, syscall.SIGTERM)
	<-a.quit
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
	if !a.closed {
		a.closed = true
		for idx := range a.closers {
			a.closers[idx]()
		}
		a.quit <- os.Kill
	}
	a.closeLock.Unlock()
}

func (a *app) traceRequest(
	context context.Context,
	method,
	path,
	query,
	agent,
	ip string,
	status,
	bytes int,
	start,
	end time.Time,
	ingress string,
) {
	latency := end.Sub(start)

	lgr := a.lgrf.Create(context)
	a.trc.TraceRequest(
		context,
		method,
		path,
		query,
		status,
		bytes,
		ip,
		agent,
		start,
		end,
		map[string]string{
			"ingress": ingress,
		},
	)
	lgr.Info(
		"Request",
		zap.Int("status", status),
		zap.String("method", method),
		zap.String("path", path),
		zap.String("query", query),
		zap.String("ip", ip),
		zap.String("userAgent", agent),
		zap.Time("mvts", end),
		zap.String("pmvts", end.Format("2006-01-02T15:04:05-0700")),
		zap.Duration("latency", latency),
		zap.String("pLatency", latency.String()),
	)
}
