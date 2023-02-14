package repos

import (
	"context"
	"prototodo/pkg/domain/base/logger"
	"prototodo/pkg/domain/domains/quotes"

	"github.com/BetaLixT/tsqlx"
	"github.com/betalixt/gorr"
)

type QuotesRepository struct {
	BaseDataRepository
	lgrf logger.IFactory
}

func NewQuotesRepository(
	dbctx *tsqlx.TracedDB,
	lgrf logger.IFactory,
) *QuotesRepository {
	return &QuotesRepository{
		BaseDataRepository: BaseDataRepository{
			dbctx: dbctx,
		},
		lgrf: lgrf,
	}
}

var _ quotes.IRepository = (*QuotesRepository)(nil)

func (r *QuotesRepository) Create(
	ctx context.Context,
	id string,
	quote string,
) (*quotes.QuoteEvent, error) {
	return nil, gorr.NewNotImplemented()
}

func (r *QuotesRepository) GetRandom(
	ctx context.Context,
) (*quotes.Quote, error) {
	return nil, gorr.NewNotImplemented()
}
