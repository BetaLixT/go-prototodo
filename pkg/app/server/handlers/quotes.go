package handlers

import (
	"context"
	srvcontracts "prototodo/pkg/app/server/contracts"
	"prototodo/pkg/domain/base/cntxt"
	"prototodo/pkg/domain/base/logger"
	"prototodo/pkg/domain/contracts"
	"time"

	"go.uber.org/zap"
)

type QuotesHandler struct {
	srvcontracts.UnimplementedQuotesServer
	ctxf cntxt.IFactory
	lgrf logger.IFactory
}

var _ srvcontracts.QuotesServer = (*QuotesHandler)(nil)

func NewQuotesHandler(
	ctxf cntxt.IFactory,
	lgrf logger.IFactory,
) *QuotesHandler {
	return &QuotesHandler{
		ctxf: ctxf,
		lgrf: lgrf,
	}
}
func (h *QuotesHandler) Get(
	c context.Context,
	qry *contracts.GetQuoteQuery,
) (res *contracts.QuoteData, err error) {
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
