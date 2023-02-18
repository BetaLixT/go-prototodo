package repos

import (
	"context"
	"prototodo/pkg/domain/base/logger"
	"prototodo/pkg/domain/base/uniques"
	"prototodo/pkg/infra/impls/evcqrs/cntxt"
	"prototodo/pkg/infra/impls/evcqrs/common"
	"prototodo/pkg/infra/impls/evcqrs/entities"

	"github.com/BetaLixT/tsqlx"
	"go.uber.org/zap"
)

type UniqueRepository struct {
	*BaseDataRepository
	lgrf logger.IFactory
}

var _ uniques.IRepository = (*UniqueRepository)(nil)

func NewUniqueRepository(
	base *BaseDataRepository,
	lgrf logger.IFactory,
	db *tsqlx.TracedDB,
) *UniqueRepository {
	return &UniqueRepository{
		BaseDataRepository: base,
		lgrf:               lgrf,
	}
}

func (r *UniqueRepository) RegisterConstraint(
	c context.Context,
	stream string,
	streamId string,
	sagaId *string,
	property string,
	value string,
) error {
	lgr := r.lgrf.Create(c)

	ctx, ok := c.(cntxt.IContext)
	if !ok {
		lgr.Error("unexpected context type")
		return common.NewFailedToAssertContextTypeError()
	}

	dbtx, err := r.getDBTx(ctx)
	if err != nil {
		lgr.Error("failed to get db transaction", zap.Error(err))
		return err
	}

	unqDao := entities.Unique{}
	err = dbtx.Get(
		ctx,
		&unqDao,
		InsertConstraintQuery,
		stream,
		streamId,
		sagaId,
		property,
		value,
	)
	if err != nil {
		lgr.Error("failed to insert unique constraint",
			zap.Error(err),
		)
	}

	return err
}

func (r *UniqueRepository) RemoveConstraint(
	c context.Context,
	stream string,
	streamId string,
) error {
	lgr := r.lgrf.Create(c)

	ctx, ok := c.(cntxt.IContext)
	if !ok {
		lgr.Error("unexpected context type")
		return common.NewFailedToAssertContextTypeError()
	}

	dbtx, err := r.getDBTx(ctx)
	if err != nil {
		lgr.Error("failed to get db transaction", zap.Error(err))
		return err
	}

	// deleting constraint
	unqDao := entities.Unique{}
	err = dbtx.Get(
		ctx,
		&unqDao,
		DeleteConstraintQuery,
		stream,
		streamId,
	)
	if err != nil {
		lgr.Error("failed to delete unique constraint",
			zap.Error(err),
		)
	}

	return err
}

const (
	InsertConstraintQuery = `
	INSERT INTO uniques(
		stream,
		stream_id,
		saga_id,
		property,
		value
	) VALUES(
		$1, $2, $3, $4, $5
	) RETURNING *
	`
	DeleteConstraintQuery = `
	DELETE FROM uniques
	WHERE stream = $1 AND stream_id = $2
	RETURNING *
	`
)
