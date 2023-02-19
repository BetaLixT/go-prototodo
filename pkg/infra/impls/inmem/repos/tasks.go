package repos

import (
	"context"
	"prototodo/pkg/domain/domains/tasks"

	"github.com/betalixt/gorr"
)

// TasksRepository repository implimentation for tasks
type TasksRepository struct{}

// NewTasksRepository creates new TasksRepository
func NewTasksRepository() *TasksRepository {
	return &TasksRepository{}
}

var _ tasks.IRepository = (*TasksRepository)(nil)

// Create creates a new task
func (r *TasksRepository) Create(
	c context.Context,
	id string,
	sagaID *string,
	data tasks.TaskData,
) (*tasks.TaskEvent, error) {
	return nil, gorr.NewNotImplemented()
}

// Get fetches an exiting task
func (r *TasksRepository) Get(
	ctx context.Context,
	id string,
) (*tasks.Task, error) {
	return nil, gorr.NewNotImplemented()
}

// List gives a paged list of tasks
func (r *TasksRepository) List(
	ctx context.Context,
	countPerPage int,
	pageNumber int,
) ([]tasks.Task, error) {
	return nil, gorr.NewNotImplemented()
}

// Delete deletes an existing task
func (r *TasksRepository) Delete(
	c context.Context,
	id string,
	sagaID *string,
	version uint64,
) (*tasks.TaskEvent, error) {
	return nil, gorr.NewNotImplemented()
}

// Update updates an existing task
func (r *TasksRepository) Update(
	c context.Context,
	id string,
	sagaID *string,
	version uint64,
	dat tasks.TaskData,
) (*tasks.TaskEvent, error) {
	return nil, gorr.NewNotImplemented()
}
