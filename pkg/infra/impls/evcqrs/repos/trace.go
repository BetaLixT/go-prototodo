package repos

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"prototodo/pkg/domain/base/logger"
	"prototodo/pkg/domain/base/trace"
	"prototodo/pkg/infra/impls/evcqrs/common"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

var _ trace.IRepository = (*TraceRepository)(nil)

// TraceRepository repository for generating, injecting and extracting trace
// info
type TraceRepository struct {
	lgrf logger.IFactory
}

// NewTraceRepository constructs a TraceRepository
func NewTraceRepository(
	lgrf logger.IFactory,
) *TraceRepository {
	return &TraceRepository{
		lgrf: lgrf,
	}
}

// ParseTraceParent parses and or generates trace information and returns
// context with trace information injected
func (r *TraceRepository) ParseTraceParent(
	parent context.Context,
	traceprnt string,
) (context.Context, error) {
	lgr := r.lgrf.Create(parent)
	ver, tid, pid, flg, err := decodeTraceparent(traceprnt)
	// If the header could not be decoded, generate a new header
	if err != nil {
		ver, flg = "00", "01"
		if tid, err = generateRadomHexString(16); err != nil {
			lgr.Error("failed to generate trace id", zap.Error(err))
			return nil, common.NewHexStringGenerationFailedError(err)
		}
	}

	// Generate a new resource id
	rid, err := generateRadomHexString(8)
	if err != nil {
		lgr.Error("failed to generate request id", zap.Error(err))
		return nil, common.NewHexStringGenerationFailedError(err)
	}

	// Generate a transaction context usin the factory
	trc := trace.TxModel{
		Ver: ver,
		Tid: tid,
		Pid: pid,
		Rid: rid,
		Flg: flg,
	}

	return context.WithValue(parent, common.TraceKey, trc), nil
}

// ExtractTraceParent extracts injected trace information from context
func (*TraceRepository) ExtractTraceParent(
	ctx context.Context,
) trace.TxModel {
	val := ctx.Value(common.TraceKey)
	if val != nil {
		if val, ok := val.(trace.TxModel); ok {
			return val
		}
	}
	return trace.TxModel{}
}

func generateRadomHexString(n int) (string, error) {
	buff := make([]byte, n)
	if _, err := rand.Read(buff); err != nil {
		return "", err
	}
	return hex.EncodeToString(buff), nil
}

func decodeTraceparent(traceparent string) (string, string, string, string, error) {
	// Fast fail for common case of empty string
	if traceparent == "" {
		return "", "", "", "", fmt.Errorf("traceparent is empty string")
	}

	hexfmt, err := regexp.Compile("^[0-9A-Fa-f]*$")
	vals := strings.Split(traceparent, "-")

	if len(vals) == 4 {
		ver, tid, pid, flg := vals[0], vals[1], vals[2], vals[3]
		if !hexfmt.MatchString(ver) || len(ver) != 2 {
			err = fmt.Errorf("invalid traceparent version")
		} else if !hexfmt.MatchString(pid) || len(pid) != 16 {
			err = fmt.Errorf("invalid traceparent parent id")
		} else if !hexfmt.MatchString(flg) || len(flg) != 2 {
			err = fmt.Errorf("invalid traceparent flag")
		} else if !hexfmt.MatchString(tid) || len(tid) != 32 {
			err = fmt.Errorf("invalid traceparent trace id")
		} else if tid == "00000000000000000000000000000000" {
			err = fmt.Errorf("traceparent trace id value is zero")
		} else {
			return ver, tid, pid, flg, nil
		}
	} else {
		err = fmt.Errorf("invalid traceparent trace id")
	}

	return "", "", "", "", err
}
