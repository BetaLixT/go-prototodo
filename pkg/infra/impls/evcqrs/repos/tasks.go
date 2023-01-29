package repos

import (
	"context"
	"prototodo/pkg/domain/base/logger"
	"prototodo/pkg/domain/domains/tasks"
	"prototodo/pkg/infra/impls/evcqrs/cntxt"
	"prototodo/pkg/infra/impls/evcqrs/common"

	"github.com/BetaLixT/tsqlx"
	"go.uber.org/zap"
)

type TasksRepostory struct {
	dbctx tsqlx.TracedDB
	lgrf  logger.IFactory
}

var _ tasks.IRepository = (*TasksRepostory)(nil)

func (r *TasksRepostory) Create(
	c context.Context,
	id string,
	data tasks.TaskData,
) (*tasks.TaskEvent, error) {
	lgr := r.lgrf.Create(c)

	ctx, ok := c.(cntxt.IContext)
	if !ok {
		lgr.Error("unexpected context type")
		return nil, common.NewFailedToAssertContextTypeError()
	}

	dbtx, err := r.getDBTx(ctx)
	if err != nil {
		lgr.Error("failed to get db transaction", zap.Error(err))
		return nil, err
	}
}

func (r *TasksRepostory) getDBTx(
	ctx cntxt.IContext,
) (*tsqlx.TracedTx, error) {
	idbtx, err := ctx.GetTransaction(
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
		return dbtx, nil
	}
}
