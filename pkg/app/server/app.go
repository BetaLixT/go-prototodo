package server

import (
	"context"
	"prototodo/pkg/infra/impls/evcqrs"

	"github.com/google/wire"
)

func Start(impl string) {
  var a app
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

type IInterface interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type app interface {
  start(context.Context)
}


// =============================================================================
// Event sourced CQRS Implementation
// =============================================================================


var cqrsDependencySet = wire.NewSet(
	evcqrs.DependencySet,
	newAppCQRS,
)

type appCQRS struct {
}

func newAppCQRS() *appCQRS {
	return &appCQRS{}
}

func (*appCQRS) start(ctx context.Context) {

}


// =============================================================================
// In Memory Implementation
// =============================================================================


var inMemDependencySet = wire.NewSet(
	newAppInMem,
)

type appInMem struct {
}

func newAppInMem() *appInMem {
	return &appInMem{}
}

func (*appInMem) start(ctx context.Context) {

}
