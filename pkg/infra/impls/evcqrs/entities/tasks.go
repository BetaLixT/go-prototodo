package entities

import (
	"database/sql/driver"
	"errors"
	"prototodo/pkg/domain/domains/tasks"
	"time"

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
	*t = TaskData{
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

func (t *TaskData) ToDTO() *tasks.TaskData {
	return &tasks.TaskData{
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		RandomMap:   t.RandomMap,
		Metadata:    t.Metadata.AsMap(),
	}
}

func (*TaskData) ToDTOSlice(
	lhs []TaskData,
) []tasks.TaskData {

	dtos := make([]tasks.TaskData, len(lhs))
	var t *tasks.TaskData
	for idx := range lhs {
		t = lhs[idx].ToDTO()
		dtos[idx] = *t
	}
	return dtos
}

type TaskEvent struct {
	BaseEvent
	Data TaskData `db:"data"`
}

func (dao *TaskEvent) ToDTO() *tasks.TaskEvent {
	return &tasks.TaskEvent{
		EventEntity: *dao.BaseEvent.ToDTO(),
		Data:        *dao.Data.ToDTO(),
	}
}

func (*TaskEvent) ToDTOSlice(
	daos []TaskEvent,
) ([]tasks.TaskEvent, error) {
	dtos := make([]tasks.TaskEvent, len(daos))
	var temp *tasks.TaskEvent
	for idx := range daos {
		temp = daos[idx].ToDTO()
		dtos[idx] = *temp
	}
	return dtos, nil
}

type TaskReadModel struct {
	Id              string        `db:"id"`
	Title           string        `db:"title"`
	Description     string        `db:"description"`
	Status          string        `db:"status"`
	RandomMap       JsonMapString `db:"random_map"`
	Metadata        JsonObj       `db:"metadata"`
	Version         uint64        `db:"version"`
	DateTimeCreated time.Time     `db:"date_time_created"`
	DateTimeUpdated time.Time     `db:"date_time_updated"`
}

func (dao *TaskReadModel) ToDTO() (*tasks.Task, error) {
	return &tasks.Task{
		Id:              dao.Id,
		Title:           dao.Title,
		Description:     dao.Description,
		Status:          dao.Status,
		RandomMap:       dao.RandomMap,
		Metadata:        dao.Metadata,
		Version:         dao.Version,
		DateTimeUpdated: dao.DateTimeUpdated,
		DateTimeCreated: dao.DateTimeCreated,
	}, nil
}

func (*TaskReadModel) ToDTOSlice(
	daos []TaskReadModel,
) ([]tasks.Task, error) {
	dtos := make([]tasks.Task, len(daos))
	var temp *tasks.Task
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
