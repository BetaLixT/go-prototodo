package repos

import (
	"context"
	"prototodo/pkg/infra/impls/evcqrs/cntxt"
	"prototodo/pkg/infra/impls/evcqrs/common"

	"github.com/BetaLixT/tsqlx"
)

type BaseRepository struct {
	dbctx tsqlx.TracedDB
}

func (r *BaseRepository) getDBTx(
	ctx cntxt.IContext,
) (*tsqlx.TracedTx, error) {
	idbtx, nw, err := ctx.GetTransactionObject(
		common.SqlTransactionObjectKey,
		func() (interface{}, error) {
			return r.dbctx.Beginx()
		},
	)
	if err != nil {
		return nil, err
	}

	if dbtx, ok := idbtx.(*tsqlx.TracedTx); !ok {
		return nil, common.NewFailedToAssertDatabaseCtxTypeError()
	} else {
		if nw {
			ctx.RegisterCommitAction(func(ctx context.Context) error {
				return dbtx.Commit()
			})
			ctx.RegisterCompensatoryAction(func(ctx context.Context) error {
				return dbtx.Rollback()
			})
		}
		return dbtx, nil
	}
}
