package quotes

import "context"

type IRepository interface {
	GetRandom(ctx context.Context) (Quote, error)
	Create(
		ctx context.Context,
		id string,
		quote string,
	) (QuoteEvent, error)
}
