package repos

import (
	"context"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/uniques"

	"github.com/betalixt/gorr"
)

type UniquesRepository struct{}

var _ uniques.IRepository = (*UniquesRepository)(nil)

func NewUniquesRepository() *UniquesRepository {
	return &UniquesRepository{}
}

func (r *UniquesRepository) RegisterConstraint(
	c context.Context,
	stream string,
	streamId string,
	sagaId *string,
	property string,
	value string,
) error {
	return gorr.NewNotImplemented()
}

func (r *UniquesRepository) RemoveConstraint(
	c context.Context,
	stream string,
	streamId string,
) error {
	return gorr.NewNotImplemented()
}
