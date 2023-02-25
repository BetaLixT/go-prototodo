package handlers

import (
	"context"
	"fmt"
	"techunicorn.com/udc-core/prototodo/pkg/app/server/common"
	appcontr "techunicorn.com/udc-core/prototodo/pkg/app/server/contracts"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/cntxt"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/logger"
	"techunicorn.com/udc-core/prototodo/pkg/domain/contracts"
	"techunicorn.com/udc-core/prototodo/pkg/domain/domains/tasks"
	"time"

	"github.com/betalixt/gorr"
	"go.uber.org/zap"
)

var _ appcontr.TasksServer = (*TasksHandler)(nil)

type TasksHandler struct {
	appcontr.UnimplementedTasksServer
	lgrf logger.IFactory
	svc  *tasks.Service
}

func NewTasksHandler(
	lgrf logger.IFactory,
	svc *tasks.Service,
) *TasksHandler {
	return &TasksHandler{
		lgrf: lgrf,
		svc:  svc,
	}
}

func (h *TasksHandler) Create(
	c context.Context,
	cmd *contracts.CreateTaskCommand,
) (res *contracts.TaskEvent, err error) {
	if cmd.UserContext == nil {
		return nil, common.NewUserContextMissingError()
	}
	ctx, ok := c.(cntxt.IContext)
	if !ok {
		return nil, common.NewInvalidContextProvidedToHandlerError()
	}
	ctx.SetTimeout(2 * time.Minute)
	lgr := h.lgrf.Create(ctx)
	lgr.Info(
		"handling",
		zap.Any("cmd", cmd),
	)
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = gorr.NewUnexpectedError(fmt.Errorf("%v", r))
				lgr.Error(
					"root panic recovered handling request",
					zap.Any("panic", r),
					zap.Stack("stack"),
				)
			} else {
				lgr.Error(
					"root panic recovered handling request",
					zap.Error(err),
					zap.Stack("stack"),
				)
			}
			ctx.RollbackTransaction()
			ctx.Cancel()
		}
		if err != nil {
			if _, ok := err.(*gorr.Error); !ok {
				err = gorr.NewUnexpectedError(err)
			}
		}
		return
	}()
	res, err = h.svc.CreateTask(
		ctx,
		cmd,
	)
	if err != nil {
		lgr.Error(
			"command handling failed",
			zap.Error(err),
		)
		ctx.RollbackTransaction()
	} else {
		err = ctx.CommitTransaction()
		if err != nil {
			lgr.Error(
				"failed to commit transaction",
				zap.Error(err),
			)
			ctx.RollbackTransaction()
		}
	}
	ctx.Cancel()
	return
}

func (h *TasksHandler) Delete(
	c context.Context,
	cmd *contracts.DeleteTaskCommand,
) (res *contracts.TaskEvent, err error) {
	if cmd.UserContext == nil {
		return nil, common.NewUserContextMissingError()
	}
	ctx, ok := c.(cntxt.IContext)
	if !ok {
		return nil, common.NewInvalidContextProvidedToHandlerError()
	}
	ctx.SetTimeout(2 * time.Minute)
	lgr := h.lgrf.Create(ctx)
	lgr.Info(
		"handling",
		zap.Any("cmd", cmd),
	)
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = gorr.NewUnexpectedError(fmt.Errorf("%v", r))
				lgr.Error(
					"root panic recovered handling request",
					zap.Any("panic", r),
					zap.Stack("stack"),
				)
			} else {
				lgr.Error(
					"root panic recovered handling request",
					zap.Error(err),
					zap.Stack("stack"),
				)
			}
			ctx.RollbackTransaction()
			ctx.Cancel()
		}
		if err != nil {
			if _, ok := err.(*gorr.Error); !ok {
				err = gorr.NewUnexpectedError(err)
			}
		}
		return
	}()
	res, err = h.svc.DeleteTask(
		ctx,
		cmd,
	)
	if err != nil {
		lgr.Error(
			"command handling failed",
			zap.Error(err),
		)
		ctx.RollbackTransaction()
	} else {
		err = ctx.CommitTransaction()
		if err != nil {
			lgr.Error(
				"failed to commit transaction",
				zap.Error(err),
			)
			ctx.RollbackTransaction()
		}
	}
	ctx.Cancel()
	return
}

func (h *TasksHandler) Update(
	c context.Context,
	cmd *contracts.UpdateTaskCommand,
) (res *contracts.TaskEvent, err error) {
	if cmd.UserContext == nil {
		return nil, common.NewUserContextMissingError()
	}
	ctx, ok := c.(cntxt.IContext)
	if !ok {
		return nil, common.NewInvalidContextProvidedToHandlerError()
	}
	ctx.SetTimeout(2 * time.Minute)
	lgr := h.lgrf.Create(ctx)
	lgr.Info(
		"handling",
		zap.Any("cmd", cmd),
	)
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = gorr.NewUnexpectedError(fmt.Errorf("%v", r))
				lgr.Error(
					"root panic recovered handling request",
					zap.Any("panic", r),
					zap.Stack("stack"),
				)
			} else {
				lgr.Error(
					"root panic recovered handling request",
					zap.Error(err),
					zap.Stack("stack"),
				)
			}
			ctx.RollbackTransaction()
			ctx.Cancel()
		}
		if err != nil {
			if _, ok := err.(*gorr.Error); !ok {
				err = gorr.NewUnexpectedError(err)
			}
		}
		return
	}()
	res, err = h.svc.UpdateTask(
		ctx,
		cmd,
	)
	if err != nil {
		lgr.Error(
			"command handling failed",
			zap.Error(err),
		)
		ctx.RollbackTransaction()
	} else {
		err = ctx.CommitTransaction()
		if err != nil {
			lgr.Error(
				"failed to commit transaction",
				zap.Error(err),
			)
			ctx.RollbackTransaction()
		}
	}
	ctx.Cancel()
	return
}

func (h *TasksHandler) Progress(
	c context.Context,
	cmd *contracts.ProgressTaskCommand,
) (res *contracts.TaskEvent, err error) {
	if cmd.UserContext == nil {
		return nil, common.NewUserContextMissingError()
	}
	ctx, ok := c.(cntxt.IContext)
	if !ok {
		return nil, common.NewInvalidContextProvidedToHandlerError()
	}
	ctx.SetTimeout(2 * time.Minute)
	lgr := h.lgrf.Create(ctx)
	lgr.Info(
		"handling",
		zap.Any("cmd", cmd),
	)
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = gorr.NewUnexpectedError(fmt.Errorf("%v", r))
				lgr.Error(
					"root panic recovered handling request",
					zap.Any("panic", r),
					zap.Stack("stack"),
				)
			} else {
				lgr.Error(
					"root panic recovered handling request",
					zap.Error(err),
					zap.Stack("stack"),
				)
			}
			ctx.RollbackTransaction()
			ctx.Cancel()
		}
		if err != nil {
			if _, ok := err.(*gorr.Error); !ok {
				err = gorr.NewUnexpectedError(err)
			}
		}
		return
	}()
	res, err = h.svc.ProgressTask(
		ctx,
		cmd,
	)
	if err != nil {
		lgr.Error(
			"command handling failed",
			zap.Error(err),
		)
		ctx.RollbackTransaction()
	} else {
		err = ctx.CommitTransaction()
		if err != nil {
			lgr.Error(
				"failed to commit transaction",
				zap.Error(err),
			)
			ctx.RollbackTransaction()
		}
	}
	ctx.Cancel()
	return
}

func (h *TasksHandler) Complete(
	c context.Context,
	cmd *contracts.CompleteTaskCommand,
) (res *contracts.TaskEvent, err error) {
	if cmd.UserContext == nil {
		return nil, common.NewUserContextMissingError()
	}
	ctx, ok := c.(cntxt.IContext)
	if !ok {
		return nil, common.NewInvalidContextProvidedToHandlerError()
	}
	ctx.SetTimeout(2 * time.Minute)
	lgr := h.lgrf.Create(ctx)
	lgr.Info(
		"handling",
		zap.Any("cmd", cmd),
	)
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = gorr.NewUnexpectedError(fmt.Errorf("%v", r))
				lgr.Error(
					"root panic recovered handling request",
					zap.Any("panic", r),
					zap.Stack("stack"),
				)
			} else {
				lgr.Error(
					"root panic recovered handling request",
					zap.Error(err),
					zap.Stack("stack"),
				)
			}
			ctx.RollbackTransaction()
			ctx.Cancel()
		}
		if err != nil {
			if _, ok := err.(*gorr.Error); !ok {
				err = gorr.NewUnexpectedError(err)
			}
		}
		return
	}()
	res, err = h.svc.CompleteTask(
		ctx,
		cmd,
	)
	if err != nil {
		lgr.Error(
			"command handling failed",
			zap.Error(err),
		)
		ctx.RollbackTransaction()
	} else {
		err = ctx.CommitTransaction()
		if err != nil {
			lgr.Error(
				"failed to commit transaction",
				zap.Error(err),
			)
			ctx.RollbackTransaction()
		}
	}
	ctx.Cancel()
	return
}

func (h *TasksHandler) ListQuery(
	c context.Context,
	qry *contracts.ListTasksQuery,
) (res *contracts.TaskEntityList, err error) {
	if qry.UserContext == nil {
		return nil, common.NewUserContextMissingError()
	}
	ctx, ok := c.(cntxt.IContext)
	if !ok {
		return nil, common.NewInvalidContextProvidedToHandlerError()
	}
	ctx.SetTimeout(2 * time.Minute)
	lgr := h.lgrf.Create(ctx)
	lgr.Info(
		"handling",
		zap.Any("qry", qry),
	)
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = gorr.NewUnexpectedError(fmt.Errorf("%v", r))
				lgr.Error(
					"root panic recovered handling request",
					zap.Any("panic", r),
					zap.Stack("stack"),
				)
			} else {
				lgr.Error(
					"root panic recovered handling request",
					zap.Error(err),
					zap.Stack("stack"),
				)
			}
			ctx.RollbackTransaction()
			ctx.Cancel()
		}
		if err != nil {
			if _, ok := err.(*gorr.Error); !ok {
				err = gorr.NewUnexpectedError(err)
			}
		}
		return
	}()
	res, err = h.svc.QueryTask(
		ctx,
		qry,
	)
	if err != nil {
		lgr.Error(
			"command handling failed",
			zap.Error(err),
		)
		ctx.RollbackTransaction()
	} else {
		err = ctx.CommitTransaction()
		if err != nil {
			lgr.Error(
				"failed to commit transaction",
				zap.Error(err),
			)
			ctx.RollbackTransaction()
		}
	}
	ctx.Cancel()
	return
}
