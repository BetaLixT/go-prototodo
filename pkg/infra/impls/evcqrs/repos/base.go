package repos

import (
	"context"
	"prototodo/pkg/infra/impls/evcqrs/cntxt"
	"prototodo/pkg/infra/impls/evcqrs/common"

	"github.com/BetaLixT/tsqlx"
)

type BaseRepository struct {
	dbctx *tsqlx.TracedDB
}

func (r *BaseRepository) insertEvent(
	ctx cntxt.IContext,
	trctx *tsqlx.TracedTx,
	out interface{},
	sagaId *string,
	stream string,
	id string,
	version uint64,
	event string,
	data interface{},
) error {
	_, tid, _, rid, _ := ctx.GetTraceInfo()
	return trctx.Get(
		ctx,
		&out,
		InsertEventQuery,
		sagaId,
		stream,
		id,
		version,
		event,
		data,
		tid,
		rid,
	)
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

func GetValueOrDefault[v comparable](in *v) (out v) {
	if in != nil {
		out = *in
	}
	return out
}

const (
	InsertEventQuery = `
	INSERT INTO events(
		saga_id,
		stream,
		stream_id,
		version,
		event,
		trace_id,
		request_id,
		data
	) VALUES(
		$1, $2, $3, $4, $5, $6
	) RETURNING *`
)
