package handlers

import (
	"context"
	"fmt"
	srvcontracts "prototodo/pkg/app/server/contracts"
	"prototodo/pkg/domain/base/cntxt"
	"prototodo/pkg/domain/base/logger"
	"prototodo/pkg/domain/contracts"
	"prototodo/pkg/domain/domains/quotes"
	"time"

	"github.com/betalixt/gorr"
	"go.uber.org/zap"
)

type QuotesHandler struct {
	srvcontracts.UnimplementedQuotesServer
	ctxf cntxt.IFactory
	lgrf logger.IFactory
	svc  *quotes.Service
}

var _ srvcontracts.QuotesServer = (*QuotesHandler)(nil)

func NewQuotesHandler(
	ctxf cntxt.IFactory,
	lgrf logger.IFactory,
	svc *quotes.Service,
) *QuotesHandler {
	return &QuotesHandler{
		ctxf: ctxf,
		lgrf: lgrf,
		svc:  svc,
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
	res, err = h.svc.GetRandomQuote(
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

func (h *QuotesHandler) Create(
	c context.Context,
	cmd *contracts.CreateQuoteCommand,
) (res *contracts.QuoteData, err error) {
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
	res, err = h.svc.CreateQuote(
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
