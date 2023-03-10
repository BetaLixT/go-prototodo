// Package handlers handles incoming quote requets
package handlers

import (
	"context"
	"fmt"
	"techunicorn.com/udc-core/prototodo/pkg/app/server/common"
	srvcontracts "techunicorn.com/udc-core/prototodo/pkg/app/server/contracts"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/cntxt"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/logger"
	"techunicorn.com/udc-core/prototodo/pkg/domain/contracts"
	"techunicorn.com/udc-core/prototodo/pkg/domain/domains/quotes"
	"time"

	"github.com/betalixt/gorr"
	"go.uber.org/zap"
)

// QuotesHandler encapsulates handlers related to the Quote Server
type QuotesHandler struct {
	srvcontracts.UnimplementedQuotesServer
	lgrf logger.IFactory
	svc  *quotes.Service
}

var _ srvcontracts.QuotesServer = (*QuotesHandler)(nil)

// NewQuotesHandler Constructs a new QuoteHandler
func NewQuotesHandler(
	lgrf logger.IFactory,
	svc *quotes.Service,
) *QuotesHandler {
	return &QuotesHandler{
		lgrf: lgrf,
		svc:  svc,
	}
}

func (h *QuotesHandler) Get(
	c context.Context,
	qry *contracts.GetQuoteQuery,
) (res *contracts.QuoteData, err error) {
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
