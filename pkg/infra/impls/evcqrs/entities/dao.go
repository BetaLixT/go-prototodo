package entities

import (
	"database/sql/driver"
	"errors"
	"prototodo/pkg/domain/base/events"
	"prototodo/pkg/domain/domains/tasks"
	"prototodo/pkg/infra/psqldb"
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type BaseEvent struct {
	Id        uint64    `db:"id"`
	SagaId    *string   `db:"saga_id"`
	Stream    string    `db:"stream"`
	StreamId  string    `db:"stream_id"`
	Event     string    `db:"event"`
	Version   uint64    `db:"version"`
	EventTime time.Time `db:"event_time"`
}

type Unique struct {
	Stream   string  `db:"stream"`
	StreamId string  `db:"stream_id"`
	SagaId   *string `db:"saga_id"`
	Property string  `db:"property"`
	Value    string  `db:"value"`
}

// - Domain data
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

func (t *TaskData) FromDto(data *tasks.TaskData) error {
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

func (t *TaskData) FromDtoSlice(
	daos []tasks.TaskData,
) ([]TaskData, error) {
	res := make([]TaskData, len(daos))
	var err error
	for idx, dao := range daos {
		err = res[idx].FromDto(&dao)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (t *TaskData) ToDto() (*tasks.TaskData, error) {
	return &tasks.TaskData{
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		RandomMap:   t.RandomMap,
		Metadata:    t.Metadata.AsMap(),
	}, nil
}

func (*TaskData) ToDtoSlice(
	lhs []TaskData,
) ([]tasks.TaskData, error) {

	dtos := make([]tasks.TaskData, len(lhs))
	var t *tasks.TaskData
	var err error
	for idx := range lhs {
		t, err = lhs[idx].ToDto()
		dtos[idx] = *t
		if err != nil {
			return nil, err
		}
	}
	return dtos, nil
}

// - Domain events
type TaskEvent struct {
	BaseEvent
	Data TaskData `db:"data"`
}

// Mapping Functions
func (dao *BaseEvent) ToDto() (*events.EventEntity, error) {
	return &events.EventEntity{
		Id:        dao.Id,
		SagaId:    dao.SagaId,
		Stream:    dao.Stream,
		StreamId:  dao.StreamId,
		Event:     dao.Event,
		Version:   dao.Version,
		EventTime: dao.EventTime,
	}, nil
}

func (dao *TaskEvent) ToDto() (*tasks.TaskEvent, error) {

	evnt, err := dao.BaseEvent.ToDto()
	if err != nil {
		return nil, err
	}
	data, err := dao.Data.ToDto()
	if err != nil {
		return nil, err
	}

	return &tasks.TaskEvent{
		EventEntity: *evnt,
		Data:        *data,
	}, nil
}

func (*TaskEvent) ToDtoSlice(
	daos []TaskEvent,
) ([]tasks.TaskEvent, error) {
	dtos := make([]tasks.TaskEvent, len(daos))
	var temp *tasks.TaskEvent
	var err error
	for idx := range daos {
		temp, err = daos[idx].ToDto()
		if err != nil {
			return nil, err
		}
		dtos[idx] = *temp
	}
	return dtos, nil
}

// Migration script
func GetMigrationScripts() []psqldb.MigrationScript {
	migrationScripts := []psqldb.MigrationScript{
		{
			Key: "initial-event-source",
			Up: `
				CREATE TABLE events (
					id bigserial PRIMARY KEY NOT NULL,
					saga_id text,
					stream text NOT NULL,
					stream_id text NOT NULL,
					version bigint NOT NULL,
					event text NOT NULL,
					event_time timestamp with time zone NOT NULL,
					data bytea NOT NULL,
					CONSTRAINT source_unique UNIQUE (stream, stream_id, version)
				);

				CREATE INDEX idx_events_stream_events ON events(stream, stream_id);

				CREATE TRIGGER set_events_event_time
				BEFORE INSERT ON events
				FOR EACH ROW
				EXECUTE PROCEDURE trigger_set_event_time();

				CREATE TABLE uniques (
					stream text NOT NULL,
					stream_id text NOT NULL,
					saga_id text,
					property text NOT NULL,
					value text NOT NULL,
					CONSTRAINT uniques_unique_constraint UNIQUE (stream, property, value)
				);
				CREATE INDEX idx_uniques_stream_constraints ON uniques(stream, stream_id);

				CREATE TABLE acl (
					stream text NOT NULL,
					stream_id text NOT NULL,
					user_type text NOT NULL,
					user_id text NOT NULL,
					permissions int NOT NULL,
					CONSTRAINT acl_unique_constraint UNIQUE ON acl(stream, stream_id, user_type, user)
				);

				CREATE TABLE foreigns (
					stream text NOT NULL,
					stream_id text NOT NULL,
					saga_id text,
					CONSTRAINT foreigns_unique_constraint UNIQUE (stream, stream_id)
				);
				CREATE INDEX idx_foreigns_stream_constraints ON foreigns(stream, stream_id);

				CREATE TABLE foreign_constraints (
					foreign_stream text NOT NULL,
					foreign_stream_id text NOT NULL,
					stream text NOT NULL,
					stream_id NOT NULL,
					saga_id text,
					CONSTRAINT foreigns_unique_constraint UNIQUE (stream, stream_id)
				);
				`,
			Down: `
			  DROP INDEX idx_foreigns_stream_constraints;
			  DROP TABLE foreigns;
			  DROP TABLE acl;
			  DROP INDEX idx_uniques_stream_constraints;
			  DROP TABLE uniques;
				DROP TRIGGER set_events_event_time on events;
				DROP INDEX idx_events_stream_events;
				DROP TABLE events;
				`,
		},
	}
	return migrationScripts
}
