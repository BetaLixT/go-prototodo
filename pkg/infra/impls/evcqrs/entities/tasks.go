package entities

import (
	"database/sql/driver"
	"errors"
	"prototodo/pkg/domain/domains/tasks"

	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

var _ driver.Value = (*TaskData)(nil)

func (a *TaskData) Value() (driver.Value, error) {
	return proto.Marshal(a)
}

func (a *TaskData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return proto.Unmarshal(b, a)
}

func (t *TaskData) FromDTO(data *tasks.TaskData) error {
	mdata, err := structpb.NewStruct(data.Metadata)
	if err != nil {
		return err
	}
	t = &TaskData{
		Title:       data.Title,
		Description: data.Description,
		Status:      data.Status,
		RandomMap:   data.RandomMap,
		Metadata:    mdata,
	}
	return nil
}

func (t *TaskData) FromDTOSlice(
	daos []tasks.TaskData,
) ([]TaskData, error) {
	res := make([]TaskData, len(daos))
	var err error
	for idx, dao := range daos {
		err = res[idx].FromDTO(&dao)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (t *TaskData) ToDTO() (*tasks.TaskData, error) {
	return &tasks.TaskData{
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		RandomMap:   t.RandomMap,
		Metadata:    t.Metadata.AsMap(),
	}, nil
}

func (*TaskData) ToDTOSlice(
	lhs []TaskData,
) ([]tasks.TaskData, error) {

	dtos := make([]tasks.TaskData, len(lhs))
	var t *tasks.TaskData
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

type TaskEvent struct {
	BaseEvent
	Data TaskData `db:"data"`
}

func (dao *TaskEvent) ToDTO() (*tasks.TaskEvent, error) {

	evnt, err := dao.BaseEvent.ToDTO()
	if err != nil {
		return nil, err
	}
	data, err := dao.Data.ToDTO()
	if err != nil {
		return nil, err
	}

	return &tasks.TaskEvent{
		EventEntity: *evnt,
		Data:        *data,
	}, nil
}

func (*TaskEvent) ToDTOSlice(
	daos []TaskEvent,
) ([]tasks.TaskEvent, error) {
	dtos := make([]tasks.TaskEvent, len(daos))
	var temp *tasks.TaskEvent
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
