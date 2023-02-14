package entities

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"prototodo/pkg/domain/base/events"
	"prototodo/pkg/infra/psqldb"
	"time"

	_ "github.com/lib/pq"
)

type BaseEvent struct {
	Id        uint64    `db:"id"`
	SagaId    *string   `db:"saga_id"`
	Stream    string    `db:"stream"`
	StreamId  string    `db:"stream_id"`
	Event     string    `db:"event"`
	Version   uint64    `db:"version"`
	EventTime time.Time `db:"event_time"`
	TraceId   string    `db:"trace_id"`
	RequestId string    `db:"request_id"`
}

func (dao *BaseEvent) ToDTO() (*events.EventEntity) {
	return &events.EventEntity{
		Id:        dao.Id,
		SagaId:    dao.SagaId,
		Stream:    dao.Stream,
		StreamId:  dao.StreamId,
		Event:     dao.Event,
		Version:   dao.Version,
		EventTime: dao.EventTime,
	}
}

type Unique struct {
	Stream   string  `db:"stream"`
	StreamId string  `db:"stream_id"`
	SagaId   *string `db:"saga_id"`
	Property string  `db:"property"`
	Value    string  `db:"value"`
}

type JsonObj map[string]interface{}

var _ driver.Value = (*JsonObj)(nil)

func (a JsonObj) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *JsonObj) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

type JsonMapString map[string]string

var _ driver.Value = (*JsonObj)(nil)

func (a JsonMapString) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *JsonMapString) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
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
					trace_id text NOT NULL,
					request_id text NOT NULL,
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
					PRIMARY KEY(stream, stream_id)
				);
				
				CREATE TABLE acl (
					stream text NOT NULL,
					stream_id text NOT NULL,
					user_type text NOT NULL,
					user_id text NOT NULL,
					permissions int NOT NULL,	
					PRIMARY KEY(stream, stream_id, user_type, user_id)
				);

				CREATE TABLE foreigns (
					stream text NOT NULL,
					stream_id text NOT NULL,
					saga_id text,
					PRIMARY KEY(stream, stream_id)
				);

				CREATE TABLE foreign_constraints (
					foreign_stream text NOT NULL,
					foreign_stream_id text NOT NULL,
					stream text NOT NULL,
					stream_id text NOT NULL,
					saga_id text,
					PRIMARY KEY (foreign_stream, foreign_stream_id, stream, stream_id),
					CONSTRAINT foreign_constraints_fk FOREIGN KEY (foreign_stream, foreign_stream_id) REFERENCES foreigns (stream, stream_id)
				);
				`,
			Down: `
			  DROP TABLE foreign_constraints;
			  DROP TABLE foreigns;
			  DROP TABLE acl;
			  DROP TABLE uniques;
				DROP TRIGGER set_events_event_time on events;
				DROP INDEX idx_events_stream_events;
				DROP TABLE events;
				`,
		},
		{
			Key: "read-models",
			Up: `
				CREATE TABLE tasks (
					id text PRIMARY KEY NOT NULL,
					title text NOT NULL,
					description text NOT NULL,
					status text NOT NULL,
					random_map jsonb NOT NULL,
					metadata jsonb NOT NULL,

					version bigint NOT NULL,	
					date_time_created timestamp with time zone NOT NULL,
					date_time_updated timestamp with time zone NOT NULL
				);

				CREATE TRIGGER set_tasks_create_time
				BEFORE INSERT ON tasks
				FOR EACH ROW
				EXECUTE PROCEDURE trigger_set_date_time_created();

				CREATE TRIGGER set_tasks_update_time
				BEFORE UPDATE ON tasks
				FOR EACH ROW
				EXECUTE PROCEDURE trigger_set_date_time_updated();

				CREATE TABLE quotes (
					id text PRIMARY KEY NOT NULL,
					quote text NOT NULL,

					version bigint NOT NULL,
					date_time_created timestamp with time zone NOT NULL,
					date_time_updated timestamp with time zone NOT NULL
				);

				CREATE TRIGGER set_quotes_create_time
				BEFORE INSERT ON quotes
				FOR EACH ROW
				EXECUTE PROCEDURE trigger_set_date_time_created();

				CREATE TRIGGER set_quotes_update_time
				BEFORE UPDATE ON quotes
				FOR EACH ROW
				EXECUTE PROCEDURE trigger_set_date_time_updated();
				`,
			Down: `
			  DROP TRIGGER set_quotes_create_time on quotes;
			  DROP TRIGGER set_quotes_update_time on quotes;
			  DROP TABLE quotes;
			  DROP TRIGGER set_tasks_create_time on tasks;
			  DROP TRIGGER set_tasks_update_time on tasks;
			  DROP TABLE tasks;
				`,
		},
	}
	return migrationScripts
}
