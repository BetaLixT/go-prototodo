package tasks

import "time"

type TaskData struct {
	Title       *string
	Description *string
	Status      *string
	RandomMap   map[string]string
	Metadata    map[string]interface{}
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
