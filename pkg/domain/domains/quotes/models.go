package quotes

import (
	"prototodo/pkg/domain/base/events"
	"time"
)

type QuoteData struct {
	Quote *string
}

type Quote struct {
	Id              string
	Quote           string
	Version         uint64
	DateTimeUpdated time.Time
	DateTimeCreated time.Time
}

type QuoteEvent struct {
	events.EventEntity
	Data QuoteData `json:"data"`
}
