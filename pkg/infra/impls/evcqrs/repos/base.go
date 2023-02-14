package repos

import (
	"context"
	"fmt"
	"prototodo/pkg/infra/impls/evcqrs/cntxt"
	"prototodo/pkg/infra/impls/evcqrs/common"
	"strings"

	"github.com/BetaLixT/tsqlx"
)

type BaseDataRepository struct {
	dbctx *tsqlx.TracedDB
}

func (r *BaseDataRepository) insertEvent(
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
		out,
		InsertEventQuery,
		sagaId,
		stream,
		id,
		version,
		event,
		tid,
		rid,
		data,
	)
}

func (r *BaseDataRepository) getDBTx(
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

func psqlSetBuilder(
	pbeg int,
	kv ...interface{},
) (setq string, values []interface{}, err error) {
	if len(kv)%2 != 0 {
		err = common.NewUnevenKeyValueCountProvidedError()
		return
	}

	cols := make([]string, len(kv)/2)
	values = make([]interface{}, len(kv)/2)
	vidx := 0

	for idx := 0; idx < len(kv); idx += 2 {
		if kv[idx+1] == nil {
			continue
		}
		col, ok := kv[idx].(string)
		if !ok {
			err = common.NewNonStringKeyProvidedError()
			values = nil
			return
		}
		cols[vidx] = fmt.Sprintf("%s = $%d", col, pbeg+vidx)
		values[vidx] = kv[idx+1]
		vidx++
	}

	if vidx == 0 {
		err = common.NewNoValuesBeingUpdatedError()
		values = nil
		return
	}
	setq = strings.Join(cols[:vidx], ",")
	values = values[:vidx]
	return
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
		$1, $2, $3, $4, $5, $6, $7, $8
	) RETURNING *`
)
