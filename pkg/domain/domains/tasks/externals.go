package tasks

import (
	"context"
)

type IRepository interface {
	Create(
		ctx context.Context,
		id string,
		data TaskData,
	) (*TaskEvent, error)
	Get(
		ctx context.Context,
		id string,
	) (*Task, error)
	Delete(
		ctx context.Context,
		id string,
		version int,
	) (*TaskEvent, error)
	Update(
		ctx context.Context,
		id string,
		version int,
		data TaskData,
	) (*TaskEvent, error)
}
