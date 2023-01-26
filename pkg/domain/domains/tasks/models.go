package tasks

import (
	"prototodo/pkg/domain/base/events"
	"prototodo/pkg/domain/common"
	"prototodo/pkg/domain/contracts"
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

type Task struct {
	Id              string
	Title           string
	Description     string
	Status          string
	CreatedBy       string
	RandomMap       map[string]string
	Metadata        map[string]interface{}
	Version         int
	DateTimeUpdated time.Time
	DateTimeCreated time.Time
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
