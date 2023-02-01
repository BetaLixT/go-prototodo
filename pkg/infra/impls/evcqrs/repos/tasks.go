package repos

import (
	"context"
	"database/sql"
	"prototodo/pkg/domain/base/logger"
	domcom "prototodo/pkg/domain/common"
	"prototodo/pkg/domain/domains/tasks"
	"prototodo/pkg/infra/impls/evcqrs/cntxt"
	"prototodo/pkg/infra/impls/evcqrs/common"
	"prototodo/pkg/infra/impls/evcqrs/entities"

	"github.com/BetaLixT/tsqlx"
	"go.uber.org/zap"
)

type TasksRepository struct {
	BaseRepository
	lgrf logger.IFactory
}

func NewTasksRepository(
	dbctx *tsqlx.TracedDB,
	lgrf logger.IFactory,
) *TasksRepository {
	return &TasksRepository{
		BaseRepository: BaseRepository{
			dbctx: dbctx,
		},
		lgrf: lgrf,
	}
}

var _ tasks.IRepository = (*TasksRepository)(nil)

func (r *TasksRepository) Create(
	c context.Context,
	id string,
	sagaId *string,
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

	var evnt entities.TaskEvent
	err = r.insertEvent(
		ctx,
		dbtx,
		&evnt,
		sagaId,
		domcom.TaskStreamName,
		id,
		0,
		domcom.EventCreated,
		data,
	)
	if err != nil {
		return nil, err
	}

	dest := entities.TaskReadModel{}
	err = dbtx.Get(
		ctx,
		&dest,
		InsertTaskReadModelQuery,
		id,
		GetValueOrDefault(data.Title),
		GetValueOrDefault(data.Description),
		GetValueOrDefault(data.Status),
		data.RandomMap,
		data.Metadata,
		evnt.Version,
		evnt.EventTime,
		evnt.EventTime,
	)
	if err != nil {
		return nil, err
	}

	return evnt.ToDTO()
}

func (r *TasksRepository) Get(
	ctx context.Context,
	id string,
) (*tasks.Task, error) {
	var task entities.TaskReadModel
	err := r.dbctx.Get(
		ctx,
		&task,
		SelectTaskByIdQuery,
		id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domcom.NewTaskMissingError()
		}
		return nil, err
	}

	return task.ToDTO()
}

func (r *TasksRepository) List(
	ctx context.Context,
	countPerPage int,
	pageNumber int,
) ([]tasks.Task, error) {
	var tasks []entities.TaskReadModel
	err := r.dbctx.Select(
		ctx,
		&tasks,
		SelectTaskByIdQuery,
		countPerPage,
		pageNumber*countPerPage,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domcom.NewTaskMissingError()
		}
		return nil, err
	}

	return ((*entities.TaskReadModel)(nil)).ToDTOSlice(tasks)
}

func (r *TasksRepository) Delete(
	c context.Context,
	id string,
	sagaId *string,
	version uint64,
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

	var evnt entities.TaskEvent
	err = r.insertEvent(
		ctx,
		dbtx,
		&evnt,
		sagaId,
		domcom.TaskStreamName,
		id,
		version,
		domcom.EventDeleted,
		tasks.TaskData{},
	)
	if err != nil {
		return nil, err
	}

	dest := entities.TaskReadModel{}
	err = dbtx.Get(
		ctx,
		&dest,
		DeleteTaskReadModelQuery,
		id,
	)
	if err != nil {
		return nil, err
	}

	return evnt.ToDTO()
}

func (r *TasksRepository) Update(
	ctx context.Context,
	id string,
	sagaId *string,
	version uint64,
	data tasks.TaskData,
) (*tasks.TaskEvent, error) {

}

// - Queries
const (
	InsertTaskReadModelQuery = `
	INSERT INTO Task (
		id,
		title,
		description,
		status,
	  random_map,
		metadata,
		version,
		date_time_created,
		date_time_updated
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9
	) RETURNING *
	`

	DeleteTaskReadModelQuery = `
	DELETE FROM Task WHERE id = $1 RETURNING *
	`

	SelectTaskByIdQuery = `
	SELECT * FROM Task WHERE id = $1
	`

	ListTasksQuery = `
	SELECT * FROM Task LIMIT $1 OFFSET $2
	`
)
