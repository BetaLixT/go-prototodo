package repos

import (
	"context"
	"prototodo/pkg/domain/base/logger"
	"prototodo/pkg/domain/domains/quotes"

	"github.com/BetaLixT/tsqlx"
)

type QuotesRepository struct {
	BaseRepository
	lgrf logger.IFactory
}

func NewQuotesRepository(
	dbctx *tsqlx.TracedDB,
	lgrf logger.IFactory,
) *QuotesRepository {
	return &QuotesRepository{
		BaseRepository: BaseRepository{
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
) (quotes.QuoteEvent, error) {

}

func (r *QuotesRepository) GetRandom(
	ctx context.Context,
) (quotes.Quote, error) {

}
