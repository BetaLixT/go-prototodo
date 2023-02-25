// Package quotes contains all business logic and validations and DTOs around the
// quotes domain
package quotes

import (
	"context"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/logger"
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/uids"
	"techunicorn.com/udc-core/prototodo/pkg/domain/contracts"

	"go.uber.org/zap"
)

// Service encapsulates business logic and use cases around the quotes domain
type Service struct {
	repo  IRepository
	lgrf  logger.IFactory
	urepo uids.IRepository
}

// NewService constructs a new quotes Service
func NewService(
	repo IRepository,
	lgrf logger.IFactory,
	urepo uids.IRepository,
) *Service {
	return &Service{
		repo:  repo,
		lgrf:  lgrf,
		urepo: urepo,
	}
}

// GetRandomQuote business logic and validations around fetching a random quote
func (s *Service) GetRandomQuote(
	ctx context.Context,
	qry *contracts.GetQuoteQuery,
) (res *contracts.QuoteData, err error) {
	lgr := s.lgrf.Create(ctx)
	q, err := s.repo.GetRandom(ctx)
	if err == nil {
		lgr.Error("failed to create quote", zap.Error(err))
		return nil, err
	}
	res = q.ToContract()
	return
}

// CreateQuote business logic and validations around creating a quote
func (s *Service) CreateQuote(
	ctx context.Context,
	cmd *contracts.CreateQuoteCommand,
) (res *contracts.QuoteData, err error) {
	lgr := s.lgrf.Create(ctx)
	id, err := s.urepo.GetID(ctx)
	if err != nil {
		lgr.Error("failed to get unique id")
		return
	}
	q, err := s.repo.Create(
		ctx,
		id,
		cmd.SagaId,
		cmd.Quote,
	)
	if err != nil {
		lgr.Error("failed to create quote", zap.Error(err))
		return nil, err
	}
	res = q.Data.ToContract()
	return
}
