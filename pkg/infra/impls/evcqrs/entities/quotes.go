package entities

import (
	"database/sql/driver"
	"errors"
	"prototodo/pkg/domain/domains/quotes"
	"time"

	"github.com/golang/protobuf/proto"
)

var _ driver.Value = (*QuoteData)(nil)

// Value for db write
func (t *QuoteData) Value() (driver.Value, error) {
	return proto.Marshal(t)
}

// Scan to read from db
func (t *QuoteData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return proto.Unmarshal(b, t)
}

// FromDTO to populate structure from dto
func (t *QuoteData) FromDTO(data *quotes.QuoteData) error {
	*t = QuoteData{
		Quote: data.Quote,
	}
	return nil
}

// FromDTOSlice to create a dao slice from dto slice
func (t *QuoteData) FromDTOSlice(
	daos []quotes.QuoteData,
) ([]QuoteData, error) {
	res := make([]QuoteData, len(daos))
	var err error
	for idx, dao := range daos {
		err = res[idx].FromDTO(&dao)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

// ToDTO to get the dto from dao
func (t *QuoteData) ToDTO() *quotes.QuoteData {
	return &quotes.QuoteData{
		Quote: t.Quote,
	}
}

// ToDTOSlice to get the dto slice from dao slice
func (*QuoteData) ToDTOSlice(
	lhs []QuoteData,
) ([]quotes.QuoteData, error) {
	dtos := make([]quotes.QuoteData, len(lhs))
	var t *quotes.QuoteData
	for idx := range lhs {
		t = lhs[idx].ToDTO()
		dtos[idx] = *t
	}
	return dtos, nil
}

// QuoteEvent representing quote events
type QuoteEvent struct {
	BaseEvent
	Data QuoteData `db:"data"`
}

// ToDTO gets dto from dao
func (dao *QuoteEvent) ToDTO() *quotes.QuoteEvent {
	data := dao.Data.ToDTO()

	return &quotes.QuoteEvent{
		EventEntity: *dao.BaseEvent.ToDTO(),
		Data:        *data,
	}
}

// ToDTOSlice gets dto slice from dao slice
func (*QuoteEvent) ToDTOSlice(
	daos []QuoteEvent,
) []quotes.QuoteEvent {
	dtos := make([]quotes.QuoteEvent, len(daos))
	var temp *quotes.QuoteEvent
	for idx := range daos {
		temp = daos[idx].ToDTO()
		dtos[idx] = *temp
	}
	return dtos
}

// QuoteReadModel the read model for quote data
type QuoteReadModel struct {
	ID              string    `db:"id"`
	Quote           string    `db:"quote"`
	Version         uint64    `db:"version"`
	DateTimeCreated time.Time `db:"date_time_created"`
	DateTimeUpdated time.Time `db:"date_time_updated"`
}

// ToDTO gets dto from dao
func (dao *QuoteReadModel) ToDTO() *quotes.Quote {
	return &quotes.Quote{
		Id:              dao.ID,
		Quote:           dao.Quote,
		Version:         dao.Version,
		DateTimeUpdated: dao.DateTimeUpdated,
		DateTimeCreated: dao.DateTimeCreated,
	}
}

// ToDTOSlice ToDTO gets dto slice from dao slice
func (*QuoteReadModel) ToDTOSlice(
	daos []QuoteReadModel,
) []quotes.Quote {
	dtos := make([]quotes.Quote, len(daos))
	var temp *quotes.Quote
	for idx := range daos {
		temp = daos[idx].ToDTO()
		dtos[idx] = *temp
	}
	return dtos
}
