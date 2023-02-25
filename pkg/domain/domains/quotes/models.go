package quotes

import (
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/events"
	"techunicorn.com/udc-core/prototodo/pkg/domain/contracts"
	"time"
)

type QuoteData struct {
	Quote *string
}

func (m *QuoteData) ToContract() *contracts.QuoteData {
	return &contracts.QuoteData{
		Quote: m.Quote,
	}
}

type Quote struct {
	Id              string
	Quote           string
	Version         uint64
	DateTimeUpdated time.Time
	DateTimeCreated time.Time
}

func (m *Quote) ToContract() *contracts.QuoteData {
	return &contracts.QuoteData{
		Quote: &m.Quote,
	}
}

type QuoteEvent struct {
	events.EventEntity
	Data QuoteData `json:"data"`
}

// func (t *QuoteEvent) ToContract() (*contracts.QuoteEvent, error) {
// 	dat := t.Data.ToContract()
//
// 	return &contracts.QuoteEvent{
// 		Id:        t.Id,
// 		SagaId:    t.SagaId,
// 		Stream:    t.Stream,
// 		StreamId:  t.StreamId,
// 		Event:     t.Event,
// 		Version:   t.Version,
// 		EventTime: timestamppb.New(t.EventTime),
// 		Data:      dat,
// 	}, nil
// }
