package evcqrs

import (
	"context"

	"github.com/google/wire"
)

var DependencySet = wire.NewSet()

type Implementation struct {
}

func (i *Implementation) Start(ctx context.Context) error {

}

func (i *Implementation) Stop(ctx context.Context) error {

}
