package entities

import (
	"database/sql/driver"
	"errors"
	"prototodo/pkg/domain/domains/quotes"
	"time"

	"github.com/golang/protobuf/proto"
)

var _ driver.Value = (*QuoteData)(nil)

func (a *QuoteData) Value() (driver.Value, error) {
	return proto.Marshal(a)
}

func (a *QuoteData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return proto.Unmarshal(b, a)
}

func (t *QuoteData) FromDTO(data *quotes.QuoteData) error {
	*t = QuoteData{
		Quote: data.Quote,
	}
	return nil
}

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

func (t *QuoteData) ToDTO() (*quotes.QuoteData, error) {
	return &quotes.QuoteData{
		Quote: t.Quote,
	}, nil
}

func (*QuoteData) ToDTOSlice(
	lhs []QuoteData,
) ([]quotes.QuoteData, error) {

	dtos := make([]quotes.QuoteData, len(lhs))
	var t *quotes.QuoteData
	var err error
	for idx := range lhs {
		t, err = lhs[idx].ToDTO()
		dtos[idx] = *t
		if err != nil {
			return nil, err
		}
	}
	return dtos, nil
}

type QuoteEvent struct {
	BaseEvent
	Data QuoteData `db:"data"`
}

func (dao *QuoteEvent) ToDTO() (*quotes.QuoteEvent, error) {

	data, err := dao.Data.ToDTO()
	if err != nil {
		return nil, err
	}

	return &quotes.QuoteEvent{
		EventEntity: *dao.BaseEvent.ToDTO(),
		Data:        *data,
	}, nil
}

func (*QuoteEvent) ToDTOSlice(
	daos []QuoteEvent,
) ([]quotes.QuoteEvent, error) {
	dtos := make([]quotes.QuoteEvent, len(daos))
	var temp *quotes.QuoteEvent
	var err error
	for idx := range daos {
		temp, err = daos[idx].ToDTO()
		if err != nil {
			return nil, err
		}
		dtos[idx] = *temp
	}
	return dtos, nil
}

type QuoteReadModel struct {
	Id              string    `db:"id"`
	Quote           string    `db:"quote"`
	Version         uint64    `db:"version"`
	DateTimeCreated time.Time `db:"date_time_created"`
	DateTimeUpdated time.Time `db:"date_time_updated"`
}

func (dao *QuoteReadModel) ToDTO() (*quotes.Quote, error) {
	return &quotes.Quote{
		Id:              dao.Id,
		Quote:           dao.Quote,
		Version:         dao.Version,
		DateTimeUpdated: dao.DateTimeUpdated,
		DateTimeCreated: dao.DateTimeCreated,
	}, nil
}

func (*QuoteReadModel) ToDTOSlice(
	daos []QuoteReadModel,
) ([]quotes.Quote, error) {
	dtos := make([]quotes.Quote, len(daos))
	var temp *quotes.Quote
	var err error
	for idx := range daos {
		temp, err = daos[idx].ToDTO()
		if err != nil {
			return nil, err
		}
		dtos[idx] = *temp
	}
	return dtos, nil
}
