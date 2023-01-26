package handlers

import (
	"context"
	appcontr "prototodo/pkg/app/server/contracts"
	"prototodo/pkg/domain/base"
	"prototodo/pkg/domain/contracts"
	"time"

	"go.uber.org/zap"
)

var _ appcontr.TasksServer = (*TasksHandler)(nil)

type TasksHandler struct {
	ctxf base.IContextFactory
	lgrf base.ILoggerFactory
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
	out, err = TODOReplaceWithServiceFunction(
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
