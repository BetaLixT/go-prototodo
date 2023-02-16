package entities

import (
	"database/sql/driver"
	"errors"
	"prototodo/pkg/domain/domains/tasks"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

var _ driver.Value = (*TaskData)(nil)

// Value for db write
func (t *TaskData) Value() (driver.Value, error) {
	return proto.Marshal(t)
}

// Scan to read from db
func (t *TaskData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return proto.Unmarshal(b, t)
}

// FromDTO to populate structure from dto
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

// GeneratePSQLReadModelSet Generates a psql set for the read model
func (t *TaskData) GeneratePSQLReadModelSet(
	pbeg int,
) (string, []interface{}, int) {
	sets := make([]string, 5)
	vals := make([]interface{}, 5)
	idx := 0
	if t.Title != nil {
		sets[idx] = "title = $" + strconv.Itoa(pbeg)
		vals[idx] = t.Title
		idx++
		pbeg++
	}
	if t.Description != nil {
		sets[idx] = "description = $" + strconv.Itoa(pbeg)
		vals[idx] = t.Description
		idx++
		pbeg++
	}
	if t.Status != nil {
		sets[idx] = "status = $" + strconv.Itoa(pbeg)
		vals[idx] = t.Status
		idx++
		pbeg++
	}
	if t.RandomMap != nil {
		sets[idx] = "random_map = $" + strconv.Itoa(pbeg)
		vals[idx] = JsonMapString(t.RandomMap)
		idx++
		pbeg++
	}
	if t.Metadata != nil {
		sets[idx] = "metadata = $" + strconv.Itoa(pbeg)
		vals[idx] = JsonObj(t.Metadata.AsMap())
		idx++
		pbeg++
	}
	var set string
	if idx == 0 {
		set = ""
	} else {
		set = strings.Join(sets[:idx], ",")
	}
	return set, vals[:idx], pbeg
}

// FromDTOSlice to create a dao slice from dto slice
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

// ToDTO to get the dto from dao
func (t *TaskData) ToDTO() *tasks.TaskData {
	return &tasks.TaskData{
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		RandomMap:   t.RandomMap,
		Metadata:    t.Metadata.AsMap(),
	}
}

// ToDTOSlice to get the dto slice from dao slice
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

// TaskEvent representing task events
type TaskEvent struct {
	BaseEvent
	Data TaskData `db:"data"`
}

// ToDTO gets dto from dao
func (dao *TaskEvent) ToDTO() *tasks.TaskEvent {
	return &tasks.TaskEvent{
		EventEntity: *dao.BaseEvent.ToDTO(),
		Data:        *dao.Data.ToDTO(),
	}
}

// ToDTOSlice gets dto slice from dao slice
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

// TaskReadModel the read model for task data
type TaskReadModel struct {
	ID              string        `db:"id"`
	Title           string        `db:"title"`
	Description     string        `db:"description"`
	Status          string        `db:"status"`
	RandomMap       JsonMapString `db:"random_map"`
	Metadata        JsonObj       `db:"metadata"`
	Version         uint64        `db:"version"`
	DateTimeCreated time.Time     `db:"date_time_created"`
	DateTimeUpdated time.Time     `db:"date_time_updated"`
}

// ToDTO gets dto from dao
func (dao *TaskReadModel) ToDTO() (*tasks.Task, error) {
	return &tasks.Task{
		Id:              dao.ID,
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

// ToDTOSlice ToDTO gets dto slice from dao slice
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
