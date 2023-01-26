package tasks

import (
	"prototodo/pkg/domain/base/acl"
	"prototodo/pkg/domain/base/cntxt"
	"prototodo/pkg/domain/base/logger"
	"prototodo/pkg/domain/base/uid"
	"prototodo/pkg/domain/common"
	"prototodo/pkg/domain/contracts"

	"go.uber.org/zap"
)

type TaskService struct {
	repo IRepository
	lgrf logger.IFactory
	aclr acl.IRepository
	uidr uid.IRepository
}

func (s *TaskService) CreateTask(
	ctx cntxt.IContext,
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

	err := ctx.BeginTransaction()
	if err != nil {
		lgr.Error(
			"failed to begin transaction",
			zap.Error(err),
		)
		return nil, err
	}

	id, err := s.uidr.GetId(ctx)
	if err != nil {
		lgr.Error(
			"failed to get unique id",
			zap.Error(err),
		)
		ctx.RollbackTransaction()
		return nil, err
	}

	res, err := s.repo.Create(
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
		ctx.RollbackTransaction()
		return nil, err
	}

	err = s.aclr.CreateACLEntry(
		ctx,
		common.TaskStreamName,
		id,
		cmd.UserContext.UserType,
		cmd.UserContext.Id,
		acl.Read | acl.Write,
	)
	if err == nil {
		err = ctx.CommitTransaction()
	}
	if err != nil {
		lgr.Error("error while creating task", zap.Error(err))
		ctx.RollbackTransaction()
	}

	return res, err
}

func (s *TaskService) DeleteTask(
	ctx cntxt.IContext,
	cmd *contracts.DeleteTaskCommand,
) (*contracts.TaskEvent, error) {

	lgr := s.lgrf.Create(ctx)
	lgr.Info("deleting task")
	err := ctx.BeginTransaction()
	if err != nil {
		lgr.Error(
			"failed to begin transaction",
			zap.Error(err),
		)
		return nil, err
	}

	if err == nil {
		err = ctx.CommitTransaction()
	}
	if err != nil {
		lgr.Error("error while creating task", zap.Error(err))
		ctx.RollbackTransaction()
	}

	return res, err
}
