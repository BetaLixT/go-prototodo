package tasks

import (
	"prototodo/pkg/domain/base/acl"
	"prototodo/pkg/domain/base/cntxt"
	"prototodo/pkg/domain/base/logger"
	"prototodo/pkg/domain/base/uid"
	"prototodo/pkg/domain/common"
	"prototodo/pkg/domain/contracts"

	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type TaskService struct {
	repo IRepository
	lgrf logger.IFactory
	aclr acl.IRepository
	uidr uid.IRepository
}

func (s *TaskService) CreateTask(
	ctx context.Context,
	cmd *contracts.CreateTaskCommand,
) (*contracts.TaskEvent, error) {
	lgr := s.lgrf.Create(ctx)

	lgr.Info("creating task")

	// business logic validations happen here
	if cmd.UserContext.UserType != common.UserTypeUser {
		lgr.Error("only users allowed to create task")
		return nil, common.NewInvalidUserTypeForTaskError()
	}

	pending := contracts.Status(contracts.Status_PENDING).String()

	id, err := s.uidr.GetId(ctx)
	if err != nil {
		lgr.Error(
			"failed to get unique id",
			zap.Error(err),
		)
		return nil, err
	}

	evnt, err := s.repo.Create(
		ctx,
		id,
		TaskData{
			Title:       &cmd.Title,
			Description: &cmd.Description,
			Status:      &pending,
		},
	)
	if err != nil {
		lgr.Error(
			"failed to create task",
			zap.Error(err),
		)
		return nil, err
	}

	res, err := evnt.ToContract()
	if err != nil {
		lgr.Error("failed to map to contract", zap.Error(err))
		return nil, err
	}

	err = s.aclr.CreateACLEntry(
		ctx,
		common.TaskStreamName,
		id,
		cmd.UserContext.UserType,
		cmd.UserContext.Id,
		acl.Read|acl.Write,
	)
	if err != nil {
		lgr.Error("failed to create acl entry", zap.Error(err))
		return nil, err
	}

	return res, err
}

func (s *TaskService) DeleteTask(
	ctx cntxt.IContext,
	cmd *contracts.DeleteTaskCommand,
) (*contracts.TaskEvent, error) {

	lgr := s.lgrf.Create(ctx)
	lgr.Info("deleting task")

	err := s.aclr.CanWrite(
		ctx,
		common.TaskStreamName,
		cmd.Id,
		cmd.UserContext.UserType,
		cmd.UserContext.Id,
	)
	if err != nil {
		lgr.Error(
			"failure while checking acl",
			zap.Error(err),
		)
		return nil, err
	}

	task, err := s.repo.Get(ctx, cmd.Id)
	if err != nil {
		lgr.Error(
			"failed to fetch task",
			zap.Error(err),
		)
		return nil, err
	}

	err = ctx.BeginTransaction()
	if err != nil {
		lgr.Error(
			"failed to begin transaction",
			zap.Error(err),
		)
		return nil, err
	}

	evnt, err := s.repo.Delete(
		ctx,
		cmd.Id,
		task.Version+1,
	)
	if err != nil {
		lgr.Error("failed to delete task", zap.Error(err))
		ctx.RollbackTransaction()
		return nil, err
	}

	res, err := evnt.ToContract()
	if err != nil {
		lgr.Error("failed to map to contract", zap.Error(err))
		ctx.RollbackTransaction()
		return nil, err
	}

	err = s.aclr.DeleteACLEntries(
		ctx,
		common.TaskStreamName,
		cmd.Id,
	)
	if err != nil {
		lgr.Error("failed to delete acl entries", zap.Error(err))
		ctx.RollbackTransaction()
		return nil, err
	}

	err = ctx.CommitTransaction()
	if err != nil {
		lgr.Error("failed to commit transaction", zap.Error(err))
		ctx.RollbackTransaction()
	}

	return res, err
}

func (s *TaskService) UpdateTask(
	ctx cntxt.IContext,
	cmd *contracts.UpdateTaskCommand,
) (*contracts.TaskEvent, error) {

	lgr := s.lgrf.Create(ctx)
	lgr.Info("updating task")

	err := s.aclr.CanWrite(
		ctx,
		common.TaskStreamName,
		cmd.Id,
		cmd.UserContext.UserType,
		cmd.UserContext.Id,
	)
	if err != nil {
		lgr.Error(
			"failure while checking acl",
			zap.Error(err),
		)
		return nil, err
	}

	task, err := s.repo.Get(ctx, cmd.Id)
	if err != nil {
		lgr.Error(
			"failed to fetch task",
			zap.Error(err),
		)
		return nil, err
	}

	err = ctx.BeginTransaction()
	if err != nil {
		lgr.Error(
			"failed to begin transaction",
			zap.Error(err),
		)
		return nil, err
	}

	evnt, err := s.repo.Update(
		ctx,
		cmd.Id,
		task.Version+1,
		TaskData{
			Title:       cmd.Title,
			Description: cmd.Description,
		},
	)
	if err != nil {
		lgr.Error("failed to update task", zap.Error(err))
		ctx.RollbackTransaction()
		return nil, err
	}

	res, err := evnt.ToContract()
	if err != nil {
		lgr.Error("failed to map to contract", zap.Error(err))
		ctx.RollbackTransaction()
		return nil, err
	}

	err = ctx.CommitTransaction()
	if err != nil {
		lgr.Error("failed to commit transaction", zap.Error(err))
		ctx.RollbackTransaction()
	}

	return res, err
}

func (s *TaskService) ProgressTask(
	ctx cntxt.IContext,
	cmd *contracts.ProgressTaskCommand,
) (*contracts.TaskEvent, error) {

	lgr := s.lgrf.Create(ctx)
	lgr.Info("progressing task")

	err := s.aclr.CanWrite(
		ctx,
		common.TaskStreamName,
		cmd.Id,
		cmd.UserContext.UserType,
		cmd.UserContext.Id,
	)
	if err != nil {
		lgr.Error(
			"failure while checking acl",
			zap.Error(err),
		)
		return nil, err
	}

	task, err := s.repo.Get(ctx, cmd.Id)
	if err != nil {
		lgr.Error(
			"failed to fetch task",
			zap.Error(err),
		)
		return nil, err
	}

	if task.Status != contracts.Status_PENDING.String() {
		lgr.Error(
			"can't progress task that isn't pending",
		)
		return nil, common.NewNotPendingTaskError()
	}

	err = ctx.BeginTransaction()
	if err != nil {
		lgr.Error(
			"failed to begin transaction",
			zap.Error(err),
		)
		return nil, err
	}

	progress := contracts.Status(contracts.Status_PROGRESS).String()
	evnt, err := s.repo.Update(
		ctx,
		cmd.Id,
		task.Version+1,
		TaskData{
			Status: &progress,
		},
	)
	if err != nil {
		lgr.Error("failed to update task", zap.Error(err))
		ctx.RollbackTransaction()
		return nil, err
	}

	res, err := evnt.ToContract()
	if err != nil {
		lgr.Error("failed to map to contract", zap.Error(err))
		ctx.RollbackTransaction()
		return nil, err
	}

	err = ctx.CommitTransaction()
	if err != nil {
		lgr.Error("failed to commit transaction", zap.Error(err))
		ctx.RollbackTransaction()
	}

	return res, err
}

func (s *TaskService) CompleteTask(
	ctx cntxt.IContext,
	cmd *contracts.CompleteTaskCommand,
) (*contracts.TaskEvent, error) {

	lgr := s.lgrf.Create(ctx)
	lgr.Info("completing task")

	err := s.aclr.CanWrite(
		ctx,
		common.TaskStreamName,
		cmd.Id,
		cmd.UserContext.UserType,
		cmd.UserContext.Id,
	)
	if err != nil {
		lgr.Error(
			"failure while checking acl",
			zap.Error(err),
		)
		return nil, err
	}

	task, err := s.repo.Get(ctx, cmd.Id)
	if err != nil {
		lgr.Error(
			"failed to fetch task",
			zap.Error(err),
		)
		return nil, err
	}

	if task.Status != contracts.Status_PROGRESS.String() {
		lgr.Error(
			"can't complete task that isn't in progress",
		)
		return nil, common.NewNotProgressTaskError()
	}

	err = ctx.BeginTransaction()
	if err != nil {
		lgr.Error(
			"failed to begin transaction",
			zap.Error(err),
		)
		return nil, err
	}

	completed := contracts.Status(contracts.Status_COMPLETED).String()
	evnt, err := s.repo.Update(
		ctx,
		cmd.Id,
		task.Version+1,
		TaskData{
			Status: &completed,
		},
	)
	if err != nil {
		lgr.Error("failed to update task", zap.Error(err))
		ctx.RollbackTransaction()
		return nil, err
	}

	res, err := evnt.ToContract()
	if err != nil {
		lgr.Error("failed to map to contract", zap.Error(err))
		ctx.RollbackTransaction()
		return nil, err
	}

	err = ctx.CommitTransaction()
	if err != nil {
		lgr.Error("failed to commit transaction", zap.Error(err))
		ctx.RollbackTransaction()
	}

	return res, err
}

func (s *TaskService) QueryTask(
	ctx cntxt.IContext,
	qry *contracts.ListTasksQuery,
) (*contracts.TaskEntityList, error) {

	lgr := s.lgrf.Create(ctx)
	lgr.Info("query tasks")

	tasks, err := s.repo.List(ctx, int(qry.CountPerPage), int(qry.PageNumber))
	if err != nil {
		lgr.Error(
			"failed to fetch task",
			zap.Error(err),
		)
		return nil, err
	}

	tasksctr, err := (*Task)(nil).ToContractSlice(tasks)
	if err != nil {
		lgr.Error("failed to map to contract", zap.Error(err))
		return nil, err
	}
	res := &contracts.TaskEntityList{
		Tasks: tasksctr,
	}

	return res, err
}
