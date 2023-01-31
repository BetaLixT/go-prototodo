package entities

import (
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
}

func (dao *BaseEvent) ToDTO() (*events.EventEntity, error) {
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

type Unique struct {
	Stream   string  `db:"stream"`
	StreamId string  `db:"stream_id"`
	SagaId   *string `db:"saga_id"`
	Property string  `db:"property"`
	Value    string  `db:"value"`
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
					PRIMARY KEY(stream, stream_id)
				);
				
				CREATE TABLE acl (
					stream text NOT NULL,
					stream_id text NOT NULL,
					user_type text NOT NULL,
					user_id text NOT NULL,
					permissions int NOT NULL,	
					PRIMARY KEY(stream, stream_id, user_type, user)
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
					stream_id NOT NULL,
					saga_id text,
					PRIMARY KEY (foreign_stream, foreign_stream_id, stream, stream_id)
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
	}
	return migrationScripts
}
