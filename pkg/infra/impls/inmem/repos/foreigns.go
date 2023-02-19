package repos

import (
	"context"
	"prototodo/pkg/domain/base/foreigns"

	"github.com/betalixt/gorr"
)

type ForeignsRepository struct{}

var _ foreigns.IRepository = (*ForeignsRepository)(nil)

func NewForeignsRepository() *ForeignsRepository {
	return &ForeignsRepository{}
}

func (r *ForeignsRepository) RegisterForeignItem(
	c context.Context,
	sagaId *string,
	foreignStream string,
	foreignStreamId string,
) error {
	return gorr.NewNotImplemented()
}

func (r *ForeignsRepository) RemoveForeignItem(
	c context.Context,
	foreignStream string,
	foreignStreamId string,
) error {
	return gorr.NewNotImplemented()
}

func (r *ForeignsRepository) RegisterConstraint(
	c context.Context,
	sagaId *string,
	foreignStream string,
	foreignStreamId string,
	stream string,
	streamId string,
) error {
	return gorr.NewNotImplemented()
}

func (r *ForeignsRepository) RemoveConstraint(
	c context.Context,
	foreignStream string,
	foreignStreamId string,
	stream string,
	streamId string,
) error {
	return gorr.NewNotImplemented()
}

func (r *ForeignsRepository) ListAssociatedObjects(
	c context.Context,
	foreignStream string,
	foreignStreamId string,
) ([]foreigns.Object, error) {
	return nil, gorr.NewNotImplemented()
}
