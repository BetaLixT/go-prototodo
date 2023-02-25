package tasks

import (
	"techunicorn.com/udc-core/prototodo/pkg/domain/base/events"
	"techunicorn.com/udc-core/prototodo/pkg/domain/common"
	"techunicorn.com/udc-core/prototodo/pkg/domain/contracts"
	"time"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TaskData struct {
	Title       *string
	Description *string
	Status      *string
	RandomMap   map[string]string
	Metadata    map[string]interface{}
}

func (t *TaskData) ToContract() (*contracts.TaskData, error) {
	cntr := &contracts.TaskData{
		Title:       t.Title,
		Description: t.Description,
		RandomMap:   t.RandomMap,
	}
	if t.Status != nil {
		if s, ok := contracts.Status_value[*t.Status]; ok {
			sc := contracts.Status(s)
			cntr.Status = &sc
		} else {
			return nil, common.NewInvalidTaskStatusError()
		}
	}
	if t.Metadata != nil {
		var err error
		cntr.Metadata, err = structpb.NewStruct(t.Metadata)
		if err != nil {
			return nil, err
		}
	}
	return cntr, nil
}

func (*TaskData) ToContractSlice(in []TaskData) ([]*contracts.TaskData, error) {
	res := make([]*contracts.TaskData, len(in))
	var err error
	for idx, t := range in {
		res[idx], err = t.ToContract()
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

type Task struct {
	Id              string
	Title           string
	Description     string
	Status          string
	RandomMap       map[string]string
	Metadata        map[string]interface{}
	Version         uint64
	DateTimeUpdated time.Time
	DateTimeCreated time.Time
}

func (t *Task) ToContract() (*contracts.TaskEntity, error) {
	s, ok := contracts.Status_value[t.Status]
	if !ok {
		return nil, common.NewInvalidTaskStatusError()
	}
	return &contracts.TaskEntity{
		Id:              t.Id,
		Version:         t.Version,
		Title:           t.Title,
		Description:     t.Description,
		Status:          contracts.Status(s),
		CreatedDateTime: timestamppb.New(t.DateTimeCreated),
		UpdatedDateTime: timestamppb.New(t.DateTimeUpdated),
	}, nil
}

func (*Task) ToContractSlice(in []Task) ([]*contracts.TaskEntity, error) {
	res := make([]*contracts.TaskEntity, len(in))
	var err error
	for idx, t := range in {
		res[idx], err = t.ToContract()
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

type TaskEvent struct {
	events.EventEntity
	Data TaskData `json:"data"`
}

func (t *TaskEvent) ToContract() (*contracts.TaskEvent, error) {
	dat, err := t.Data.ToContract()
	if err != nil {
		return nil, err
	}

	return &contracts.TaskEvent{
		Id:        t.Id,
		SagaId:    t.SagaId,
		Stream:    t.Stream,
		StreamId:  t.StreamId,
		Event:     t.Event,
		Version:   t.Version,
		EventTime: timestamppb.New(t.EventTime),
		Data:      dat,
	}, nil
}

func (*TaskEvent) ToContractSlice(in []TaskEvent) ([]*contracts.TaskEvent, error) {
	res := make([]*contracts.TaskEvent, len(in))
	var err error
	for idx, t := range in {
		res[idx], err = t.ToContract()
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
