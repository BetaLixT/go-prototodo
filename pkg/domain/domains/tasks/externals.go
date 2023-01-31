package tasks

import (
	"context"
)

type IRepository interface {
	Create(
		ctx context.Context,
		id string,
		sagaId *string,
		data TaskData,
	) (*TaskEvent, error)
	Get(
		ctx context.Context,
		id string,
	) (*Task, error)
	List(
		ctx context.Context,
		countPerPage int,
		pageNumber int,
	) ([]Task, error)
	Delete(
		ctx context.Context,
		id string,
		sagaId *string,
		version uint64,
	) (*TaskEvent, error)
	Update(
		ctx context.Context,
		id string,
		sagaId *string,
		version uint64,
		data TaskData,
	) (*TaskEvent, error)
}
