// Package server contains server logic to handle incoming requests and command
// query handlers
package server

import (
	"context"
	"prototodo/pkg/domain"
	"prototodo/pkg/domain/base/impl"
	"prototodo/pkg/domain/domains/quotes"
	"prototodo/pkg/domain/domains/tasks"
	"prototodo/pkg/infra/impls/evcqrs"
	"prototodo/pkg/infra/impls/inmem"

	"github.com/google/wire"
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
	newApp,
)

// inMemDependencySet dependency set with in memory implementation
var inMemDependencySet = wire.NewSet(
	inmem.DependencySet,
	domain.DependencySet,
	newApp,
)

// =============================================================================
// Application
// =============================================================================

type app struct {
	impl impl.IImplementation
	tsvc *tasks.Service
	qsvc *quotes.Service
}

func newApp(
	impl impl.IImplementation,
	tsvc *tasks.Service,
	qsvc *quotes.Service,
) *app {
	return &app{
		impl: impl,
		tsvc: tsvc,
		qsvc: qsvc,
	}
}

func (*app) start(ctx context.Context) {
}
