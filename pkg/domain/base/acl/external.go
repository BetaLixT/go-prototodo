package acl

import "context"

type IRepository interface {
	CreateACLEntry(
		ctx context.Context,
		stream string,
		streamId string,
		userType string,
		userId string,
		permissions int,
	) error
	DeleteACLEntry(ctx context.Context,
		stream string,
		streamId string,
		userType string,
		userId string,
	) error
	CanRead(ctx context.Context,
		stream string,
		streamId string,
		userType string,
		userId string,
	) error
	CanWrite(ctx context.Context,
		stream string,
		streamId string,
		userType string,
		userId string,
	) error
}
