package quotes

import (
	"context"
	"prototodo/pkg/domain/contracts"
)

type QuotesService struct {
	repo IRepository
}

func (s *QuotesService) GetRandomQuote(
	ctx context.Context,
	qry *contracts.GetQuoteQuery,
) (*contracts.QuoteData, error) {

}

func (s *QuotesService) CreateQuote(
	ctx context.Context,
	cmd *contracts.CreateQuoteCommand,
) (*contracts.QuoteData, error) {

}
