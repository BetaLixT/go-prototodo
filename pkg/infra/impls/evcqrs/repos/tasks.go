package repos

import (
	"context"
	"database/sql"
	"fmt"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/logger"
	domcom "techunicorn.com/udc-core/prototodo/pkg/domain/common"
	"techunicorn.com/udc-core/prototodo/pkg/domain/domains/tasks"
	"techunicorn.com/udc-core/prototodo/pkg/infra/impls/evcqrs/cntxt"
	"techunicorn.com/udc-core/prototodo/pkg/infra/impls/evcqrs/common"
	"techunicorn.com/udc-core/prototodo/pkg/infra/impls/evcqrs/entities"

	"go.uber.org/zap"
)

// TasksRepository repository implimentation for tasks
type TasksRepository struct {
	*BaseDataRepository
	lgrf logger.IFactory
}

// NewTasksRepository creates new TasksRepository
func NewTasksRepository(
	base *BaseDataRepository,
	lgrf logger.IFactory,
) *TasksRepository {
	return &TasksRepository{
		BaseDataRepository: base,
		lgrf:               lgrf,
	}
}

var _ tasks.IRepository = (*TasksRepository)(nil)

// Create creates a new task
func (r *TasksRepository) Create(
	c context.Context,
	id string,
	sagaID *string,
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
	var dat entities.TaskData
	dat.FromDTO(&data)
	err = r.insertEvent(
		ctx,
		dbtx,
		&evnt,
		sagaID,
		domcom.TaskStreamName,
		id,
		0,
		domcom.EventCreated,
		&dat,
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
		entities.JSONMapString(data.RandomMap),
		entities.JSONObj(data.Metadata),
		evnt.Version,
		evnt.EventTime,
		evnt.EventTime,
	)
	if err != nil {
		return nil, err
	}

	return evnt.ToDTO(), nil
}

// Get fetches an exiting task
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

// List gives a paged list of tasks
func (r *TasksRepository) List(
	ctx context.Context,
	countPerPage int,
	pageNumber int,
) ([]tasks.Task, error) {
	var tasks []entities.TaskReadModel
	err := r.dbctx.Select(
		ctx,
		&tasks,
		ListTasksQuery,
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

// Delete deletes an existing task
func (r *TasksRepository) Delete(
	c context.Context,
	id string,
	sagaID *string,
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
		sagaID,
		domcom.TaskStreamName,
		id,
		version,
		domcom.EventDeleted,
		&entities.TaskData{},
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
		version-1,
	)
	if err != nil {
		return nil, err
	}

	return evnt.ToDTO(), nil
}

// Update updates an existing task
func (r *TasksRepository) Update(
	c context.Context,
	id string,
	sagaID *string,
	version uint64,
	dat tasks.TaskData,
) (*tasks.TaskEvent, error) {
	lgr := r.lgrf.Create(c)

	ctx, ok := c.(cntxt.IContext)
	if !ok {
		lgr.Error("unexpected context type")
		return nil, common.NewFailedToAssertContextTypeError()
	}

	var data entities.TaskData
	data.FromDTO(&dat)

	set, vals, _ := data.GeneratePSQLReadModelSet(4)
	if set == "" {
		lgr.Error("no values updated")
		return nil, domcom.NewNoTaskUpdatesError()
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
		sagaID,
		domcom.TaskStreamName,
		id,
		version,
		domcom.EventUpdated,
		&data,
	)
	if err != nil {
		lgr.Error("failed to insert update event", zap.Error(err))
		return nil, err
	}

	allvals := append([]interface{}{id, version - 1, version}, vals...)
	dest := entities.TaskReadModel{}
	err = dbtx.Get(
		ctx,
		&dest,
		fmt.Sprintf(UpdateTaskQuery, set),
		allvals...,
	)
	if err != nil {
		lgr.Error("failed to update entity", zap.Error(err))
		return nil, err
	}

	return evnt.ToDTO(), nil
}

// - Queries
const (
	InsertTaskReadModelQuery = `
	INSERT INTO tasks (
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
	DELETE FROM tasks WHERE id = $1 AND version = $2 RETURNING *
	`

	SelectTaskByIdQuery = `
	SELECT * FROM tasks WHERE id = $1
	`

	ListTasksQuery = `
	SELECT * FROM tasks LIMIT $1 OFFSET $2
	`

	UpdateTaskQuery = `
	UPDATE tasks SET %s, version = $3 WHERE id = $1 AND version = $2 RETURNING *
	`
)
