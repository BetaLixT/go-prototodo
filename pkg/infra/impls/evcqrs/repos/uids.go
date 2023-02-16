package repos

import (
	"context"

	"github.com/bwmarrin/snowflake"
)

type UIDRepository struct {
	sf *snowflake.Node
}

func NewUIDRepository(
	sf *snowflake.Node,
) *UIDRepository {
	return &UIDRepository{
		sf: sf,
	}
}

func (r *UIDRepository) GetId(
	ctx context.Context,
) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	return r.sf.Generate().String(), nil
}
