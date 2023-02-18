package repos

import (
	"context"
	"prototodo/pkg/domain/base/acl"
	"prototodo/pkg/domain/base/logger"
	domcom "prototodo/pkg/domain/common"
	"prototodo/pkg/infra/impls/evcqrs/cntxt"
	"prototodo/pkg/infra/impls/evcqrs/common"
	"prototodo/pkg/infra/impls/evcqrs/entities"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

type ACLRepository struct {
	*BaseDataRepository
	rctx    *redis.Client
	mcache  *cache.Cache
	lgrf    logger.IFactory
	keySffx string
}

func NewACLRepository(
	base *BaseDataRepository,
	rctx *redis.Client,
	mcache *cache.Cache,
	lgrf logger.IFactory,
) *ACLRepository {
	return &ACLRepository{
		BaseDataRepository: base,
		rctx:               rctx,
		mcache:             mcache,
		lgrf:               lgrf,
		keySffx:            common.ACLCacheSuffix + domcom.ServiceName + ":",
	}
}

func (r *ACLRepository) CreateACLEntry(
	c context.Context,
	stream string,
	streamID string,
	userType string,
	userID string,
	permissions int,
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

	entry := entities.ACL{}
	err = dbtx.Get(
		ctx,
		&entry,
		InsertACLQuery,
		stream,
		streamID,
		userType,
		userID,
		permissions,
	)
	if err != nil {
		lgr.Error("failed to create ACL entry")
		return err
	}

	return nil
}

func (r *ACLRepository) DeleteACLEntry(
	ctx context.Context,
	stream string,
	streamID string,
	userType string,
	userID string,
) error {
}

func (r *ACLRepository) CanRead(
	ctx context.Context,
	stream string,
	streamIDs []string,
	userType string,
	userID string,
) error {
	perm, err := r.getEntry(
		ctx,
		stream,
		streamIDs,
		userType,
		userID,
	)
	if err != nil {
		return err
	}
	if (perm & acl.Read) != 0 {
		return nil
	}
	return domcom.NewUserACLCheckFailedError()
}

func (r *ACLRepository) CanWrite(
	ctx context.Context,
	stream string,
	streamIDs []string,
	userType string,
	userID string,
) error {
	perm, err := r.getEntry(
		ctx,
		stream,
		streamIDs,
		userType,
		userID,
	)
	if err != nil {
		return err
	}
	if (perm & acl.Write) != 0 {
		return nil
	}
	return domcom.NewUserACLCheckFailedError()
}

func (r *ACLRepository) getEntries(
	ctx context.Context,
	stream string,
	streamIDs []string,
	userType string,
	userID string,
) (int, error) {
	lgr := r.lgrf.Create(ctx)

	rkey := r.keySffx + userType + ":" + userID
	entries, err := r.rctx.ZRange(
		ctx,
		rkey,
		0,
		100,
	).Result()
	if err != nil {
		// only should happen if the key has been hijacked by something else
		lgr.Error("failed to range through set", zap.Error(err))
	}

	// looping through redis entries, parsing and filtering out required entries
	parsedRes := map[string]int{}
	for idx := range entries {
		split := strings.Split(entries[idx], ":")
		if len(split) != 3 {
			lgr.Warn(
				"invalid split while parsing redis cached acl entry",
				zap.String("entry", entries[idx]),
			)
			continue
		}
		if split[0] == stream {
			per, err := strconv.Atoi(split[1])
			if err != nil {
				lgr.Warn(
					"unabled to parse acl entry's permission field to int",
					zap.String("entry", entries[idx]),
				)
			}
			parsedRes[split[1]] = per
		}
	}

	perm := 0b11
	rpipe := r.rctx.Pipeline()

	// Keeping the entire LRU cache alive if it's being used, if not used for
	// more than two hours for this stream and this user, the cache will be
	// deleted
	rpipe.ExpireAt(
		ctx,
		rkey,
		time.Now().Add(2*time.Hour),
	)
	nfidx := 0
	notFound := make([]string, len(streamIDs))
	for idx := range streamIDs {
		val, ok := parsedRes[streamIDs[idx]]
		if ok {
			// ACL is implemented as an LRU cache, we increase the score as we get
			// hits and cache values with low hits will be cleaned up
			rpipe.ZIncrBy(
				ctx,
				rkey,
				1,
				generateACLSetMember(stream, streamIDs[idx], val),
			)
			perm = perm & val
		} else {
			notFound[nfidx] = streamIDs[idx]
			nfidx++
		}
		if perm == 0 {
			rpipe.Exec(ctx)
			return perm, nil
		}
	}
}

func generateACLSetMember(stream string, id string, perm int) string {
	return stream + ":" + id + ":" + strconv.Itoa(perm)
}

func (r *ACLRepository) getEntry(
	ctx context.Context,
	stream string,
	streamID string,
	userType string,
	userID string,
) (int, error) {
	lgr := r.lgrf.Create(ctx)

	entries, err := r.rctx.ZRange(
		ctx,
		r.keySffx+userType+":"+userID,
		0,
		100,
	).Result()
	if err != nil {
		// only should happen if the key has been hijacked by something else
		lgr.Error("failed to range through set", zap.Error(err))
	}

	parsedRes := map[string]int{}
	for idx := range entries {
		split := strings.Split(entries[idx], ":")
		if len(split) != 3 {
			lgr.Warn(
				"invalid split while parsing redis cached acl entry",
				zap.String("entry", entries[idx]),
			)
			continue
		}
		if split[0] == stream {
			per, err := strconv.Atoi(split[1])
			if err != nil {
				lgr.Warn(
					"unabled to parse acl entry's permission field to int",
					zap.String("entry", entries[idx]),
				)
			}
			parsedRes[split[1]] = per
		}
	}
}

// - Queries
const (
	InsertACLQuery = `
	INSERT INTO acls (
		stream,
		stream_id,
		user_type,
		user_id,
	  permissions,
	) VALUES (
		$1, $2, $3, $4, $5
	) RETURNING *
	`
)
