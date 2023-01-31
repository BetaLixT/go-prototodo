package repos

import (
	"context"
	"prototodo/pkg/domain/base/logger"
	"prototodo/pkg/domain/domains/tasks"
	"prototodo/pkg/infra/impls/evcqrs/cntxt"
	"prototodo/pkg/infra/impls/evcqrs/common"

	"go.uber.org/zap"
)

type TasksRepository struct {
	BaseRepository
	lgrf  logger.IFactory
}

var _ tasks.IRepository = (*TasksRepostory)(nil)

func (r *TasksRepository) Create(
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

// - Queries
const (

)
