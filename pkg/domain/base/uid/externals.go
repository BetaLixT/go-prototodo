package uid

import "context"

type IRepository interface {
	GetId(ctx context.Context) (string, error)
}
