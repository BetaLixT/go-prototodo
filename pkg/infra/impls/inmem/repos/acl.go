package repos

import (
	"context"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/acl"

	"github.com/betalixt/gorr"
)

type ACLRepository struct{}

var _ acl.IRepository = (*ACLRepository)(nil)

func NewACLRepository() *ACLRepository {
	return &ACLRepository{}
}

func (r *ACLRepository) CreateACLEntry(
	c context.Context,
	stream string,
	streamID string,
	userType string,
	userID string,
	permissions int,
) error {
	return gorr.NewNotImplemented()
}

func (r *ACLRepository) DeleteACLEntry(
	c context.Context,
	stream string,
	streamID string,
	userType string,
	userID string,
) error {
	return gorr.NewNotImplemented()
}

func (r *ACLRepository) DeleteACLEntries(
	c context.Context,
	stream string,
	streamID string,
) error {
	return gorr.NewNotImplemented()
}

func (r *ACLRepository) CanRead(
	ctx context.Context,
	stream string,
	streamIDs []string,
	userType string,
	userID string,
) error {
	return gorr.NewNotImplemented()
}

func (r *ACLRepository) CanWrite(
	ctx context.Context,
	stream string,
	streamIDs []string,
	userType string,
	userID string,
) error {
	return gorr.NewNotImplemented()
}
