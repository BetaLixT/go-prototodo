package evcqrs

import (
	"context"

	"github.com/betalixt/gorr"
	"github.com/google/wire"
)

var DependencySet = wire.NewSet()

type Implementation struct {
}

func (i *Implementation) Start(ctx context.Context) error {
	return gorr.NewNotImplemented()
}

func (i *Implementation) Stop(ctx context.Context) error {
	return gorr.NewNotImplemented()
}
