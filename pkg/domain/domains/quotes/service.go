// Package quotes containing all business logic and DTO related to the quotes
// domain
package quotes

import (
	"context"
	"prototodo/pkg/domain/base/logger"
	"prototodo/pkg/domain/base/uids"
	"prototodo/pkg/domain/contracts"
)

type QuotesService struct {
	repo  IRepository
	lgrf  logger.IFactory
	urepo uids.IRepository
}

func NewQuotesService(
	repo IRepository,
	lgrf logger.IFactory,
	urepo uids.IRepository,
) *QuotesService {
	return &QuotesService{
		repo:  repo,
		lgrf:  lgrf,
		urepo: urepo,
	}
}

func (s *QuotesService) GetRandomQuote(
	ctx context.Context,
	qry *contracts.GetQuoteQuery,
) (res *contracts.QuoteData, err error) {
	q, err := s.repo.GetRandom(ctx)
	if err == nil {
		res = q.ToContract()
	}
	return res, err
}

func (s *QuotesService) CreateQuote(
	ctx context.Context,
	cmd *contracts.CreateQuoteCommand,
) (res *contracts.QuoteData, err error) {
	lgr := s.lgrf.Create(ctx)
	id, err := s.urepo.GetId(ctx)
	if err != nil {
		lgr.Error("failed to get unique id")
		return
	}
	q, err := s.repo.Create(
		ctx,
		id,
		cmd.Quote,
	)
	if err != nil {
		res = q.Data.ToContract()
	}
	return
}
