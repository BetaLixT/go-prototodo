package handlers

import (
	"context"
	appcontr "prototodo/pkg/app/server/contracts"
	"prototodo/pkg/domain/base/cntxt"
	"prototodo/pkg/domain/base/logger"
	"prototodo/pkg/domain/contracts"
	"prototodo/pkg/domain/domains/tasks"
	"time"

	"go.uber.org/zap"
)

var _ appcontr.TasksServer = (*TasksHandler)(nil)

type TasksHandler struct {
	appcontr.UnimplementedTasksServer
	ctxf cntxt.IFactory
	lgrf logger.IFactory
	tsrv tasks.TaskService
}

func (a *TasksHandler) Create(
	c context.Context,
	cmd *contracts.CreateTaskCommand,
) (out *contracts.TaskEvent, err error) {
	ctx := a.ctxf.Create(
		c,
		time.Second*5,
	)
	lgr := a.lgrf.Create(ctx)
	lgr.Info(
		"handling command",
		zap.Any("cmd", cmd),
	)
	out, err = a.tsrv.CreateTask(
		ctx,
		cmd,
	)
	if err != nil {
		lgr.Error(
			"command handling failed",
			zap.Error(err),
		)
	}
	return
}

func (h *TasksHandler) Delete(
	c context.Context,
	cmd *contracts.DeleteTaskCommand,
) (res *contracts.TaskEvent, err error) {
	ctx := h.ctxf.Create(
		c,
		time.Second*5,
	)
	lgr := h.lgrf.Create(ctx)
	lgr.Info(
		"handling",
		zap.Any("cmd", cmd),
	)
	res, err = TODOReplaceWithServiceFunction(
		ctx,
		cmd,
	)
	if err != nil {
		lgr.Error(
			"handling failed",
			zap.Error(err),
		)
	}
	return
}
func (h *TasksHandler) Update(
	c context.Context,
	cmd *contracts.UpdateTaskCommand,
) (res *contracts.TaskEvent, err error) {
	ctx := h.ctxf.Create(
		c,
		time.Second*5,
	)
	lgr := h.lgrf.Create(ctx)
	lgr.Info(
		"handling",
		zap.Any("cmd", cmd),
	)
	res, err = TODOReplaceWithServiceFunction(
		ctx,
		cmd,
	)
	if err != nil {
		lgr.Error(
			"handling failed",
			zap.Error(err),
		)
	}
	return
}
func (h *TasksHandler) Progress(
	c context.Context,
	cmd *contracts.ProgressTaskCommand,
) (res *contracts.TaskEvent, err error) {
	ctx := h.ctxf.Create(
		c,
		time.Second*5,
	)
	lgr := h.lgrf.Create(ctx)
	lgr.Info(
		"handling",
		zap.Any("cmd", cmd),
	)
	res, err = TODOReplaceWithServiceFunction(
		ctx,
		cmd,
	)
	if err != nil {
		lgr.Error(
			"handling failed",
			zap.Error(err),
		)
	}
	return
}
func (h *TasksHandler) Complete(
	c context.Context,
	cmd *contracts.CompleteTaskCommand,
) (res *contracts.TaskEvent, err error) {
	ctx := h.ctxf.Create(
		c,
		time.Second*5,
	)
	lgr := h.lgrf.Create(ctx)
	lgr.Info(
		"handling",
		zap.Any("cmd", cmd),
	)
	res, err = TODOReplaceWithServiceFunction(
		ctx,
		cmd,
	)
	if err != nil {
		lgr.Error(
			"handling failed",
			zap.Error(err),
		)
	}
	return
}
func (h *TasksHandler) ListQuery(
	c context.Context,
	qry *contracts.ListTasksQuery,
) (res *contracts.TaskEntity, err error) {
	ctx := h.ctxf.Create(
		c,
		time.Second*5,
	)
	lgr := h.lgrf.Create(ctx)
	lgr.Info(
		"handling",
		zap.Any("qry", qry),
	)
	res, err = TODOReplaceWithServiceFunction(
		ctx,
		qry,
	)
	if err != nil {
		lgr.Error(
			"handling failed",
			zap.Error(err),
		)
	}
	return
}
