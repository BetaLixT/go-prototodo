package handlers

import (
	"context"
	"fmt"
	appcontr "prototodo/pkg/app/server/contracts"
	"prototodo/pkg/domain/base/cntxt"
	"prototodo/pkg/domain/base/logger"
	"prototodo/pkg/domain/contracts"
	"prototodo/pkg/domain/domains/tasks"
	"time"

	"github.com/betalixt/gorr"
	"go.uber.org/zap"
)

var _ appcontr.TasksServer = (*TasksHandler)(nil)

type TasksHandler struct {
	appcontr.UnimplementedTasksServer
	ctxf cntxt.IFactory
	lgrf logger.IFactory
	tsrv tasks.TaskService
}

func (h *TasksHandler) Create(
	c context.Context,
	cmd *contracts.CreateTaskCommand,
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
		if _, ok := err.(*gorr.Error); !ok {
			err = gorr.NewUnexpectedError(err)
		}
		return
	}()
	res, err = h.tsrv.CreateTask(
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
	ctx := h.ctxf.Create(
		c,
		time.Second*5,
	)
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
		if _, ok := err.(*gorr.Error); !ok {
			err = gorr.NewUnexpectedError(err)
		}
		return
	}()
	res, err = h.tsrv.DeleteTask(
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
	ctx := h.ctxf.Create(
		c,
		time.Second*5,
	)
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
		if _, ok := err.(*gorr.Error); !ok {
			err = gorr.NewUnexpectedError(err)
		}
		return
	}()
	res, err = h.tsrv.UpdateTask(
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
	ctx := h.ctxf.Create(
		c,
		time.Second*5,
	)
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
		if _, ok := err.(*gorr.Error); !ok {
			err = gorr.NewUnexpectedError(err)
		}
		return
	}()
	res, err = h.tsrv.ProgressTask(
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
	ctx := h.ctxf.Create(
		c,
		time.Second*5,
	)
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
		if _, ok := err.(*gorr.Error); !ok {
			err = gorr.NewUnexpectedError(err)
		}
		return
	}()
	res, err = h.tsrv.CompleteTask(
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
	ctx := h.ctxf.Create(
		c,
		time.Second*5,
	)
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
		if _, ok := err.(*gorr.Error); !ok {
			err = gorr.NewUnexpectedError(err)
		}
		return
	}()
	res, err = h.tsrv.QueryTask(
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
