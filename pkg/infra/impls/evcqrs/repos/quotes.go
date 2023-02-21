package repos

import (
	"context"
	"prototodo/pkg/domain/base/logger"
	domcom "prototodo/pkg/domain/common"
	"prototodo/pkg/domain/domains/quotes"
	"prototodo/pkg/infra/impls/evcqrs/cntxt"
	"prototodo/pkg/infra/impls/evcqrs/common"
	"prototodo/pkg/infra/impls/evcqrs/entities"

	"github.com/BetaLixT/tsqlx"
	"github.com/betalixt/gorr"
	"go.uber.org/zap"
)

// QuotesRepository repository implementation for quotes
type QuotesRepository struct {
	BaseDataRepository
	lgrf logger.IFactory
}

// NewQuotesRepository creates QuotesRepository
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

// Create a quote
func (r *QuotesRepository) Create(
	c context.Context,
	id string,
	sagaID *string,
	quote string,
) (*quotes.QuoteEvent, error) {
	lgr := r.lgrf.Create(c)

	ctx, ok := c.(cntxt.IContext)
	if !ok {
		lgr.Error("unexpected context type")
		return nil, common.NewFailedToAssertContextTypeError()
	}

	dbtx, err := r.getDBTx(ctx)
	if err != nil {
		lgr.Error("failed to get db transaction", zap.Error(err))
		return nil, err
	}

	var ev entities.QuoteEvent
	err = r.insertEvent(
		ctx,
		dbtx,
		&ev,
		sagaID,
		domcom.QuoteStreamName,
		id,
		0,
		domcom.EventCreated,
		&entities.QuoteData{
			Quote: &quote,
		},
	)
	if err != nil {
		lgr.Error("failed to insert create event", zap.Error(err))
		return nil, err
	}

	var res entities.QuoteReadModel
	err = dbtx.Get(
		ctx,
		&res,
		InsertQuoteReadModelQuery,
		id,
		ev.Data.Quote,
		ev.Version,
		ev.EventTime,
		ev.EventTime,
	)
	if err != nil {
		lgr.Error("failed to create quote read model", zap.Error(err))
		return nil, err
	}

	return ev.ToDTO(), nil
}

// GetRandom Fetch a random quote
func (r *QuotesRepository) GetRandom(
	ctx context.Context,
) (*quotes.Quote, error) {
	return nil, gorr.NewNotImplemented()
}

// - Queries
const (
	InsertQuoteReadModelQuery = `
	INSERT INTO quotes (
		id,
		quote,
		version,
		date_time_created,
		date_time_updated
	) VALUES (
		$1, $2, $3, $4, $5
	) RETURNING *
	`

	GetQuoteCountQuery = `
	SELECT COUNT(id) FROM quotes
	`

	ListQuotesQuery = `
	SELECT * FROM quotes LIMIT $1 OFFSET $2
	`
)
