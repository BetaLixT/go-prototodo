package quotes

import "context"

// IRepository interface to repository for quote data
type IRepository interface {
	GetRandom(ctx context.Context) (*Quote, error)
	Create(
		ctx context.Context,
		id string,
		sagaID *string,
		quote string,
	) (*QuoteEvent, error)
}
