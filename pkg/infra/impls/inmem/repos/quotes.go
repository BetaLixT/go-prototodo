package repos

import (
	"context"
	"techunicorn.com/udc-core/prototodo/pkg/domain/domains/quotes"

	"github.com/betalixt/gorr"
)

// QuotesRepository repository implementation for quotes
type QuotesRepository struct{}

// NewQuotesRepository creates QuotesRepository
func NewQuotesRepository() *QuotesRepository {
	return &QuotesRepository{}
}

var _ quotes.IRepository = (*QuotesRepository)(nil)

// Create a quote
func (r *QuotesRepository) Create(
	c context.Context,
	id string,
	sagaID *string,
	quote string,
) (*quotes.QuoteEvent, error) {
	return nil, gorr.NewNotImplemented()
}

// GetRandom Fetch a random quote
func (r *QuotesRepository) GetRandom(
	ctx context.Context,
) (*quotes.Quote, error) {
	return nil, gorr.NewNotImplemented()
}
