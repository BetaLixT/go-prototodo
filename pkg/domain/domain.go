// Package domain containing business logic, contracts and models
package domain

import (
	"prototodo/pkg/domain/domains/quotes"
	"prototodo/pkg/domain/domains/tasks"

	"github.com/google/wire"
)

// DependencySet dependencies provided by the domain
var DependencySet = wire.NewSet(
	tasks.NewService,
	quotes.NewService,
)
