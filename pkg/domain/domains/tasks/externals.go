package tasks

import (
	"context"
	"prototodo/pkg/domain/contracts"
)

type IRepository interface {
	Create(
		ctx context.Context,
		data TaskData,
	) (*contracts.TaskEvent, error)
	Get(
		ctx context.Context,
		id string,
	) (Task, error)
	Delete(
		ctx context.Context,
		id string,
		version int,
	) (*contracts.TaskEvent, error)
}
