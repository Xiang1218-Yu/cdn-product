package store

import (
	"context"
	"time"

	"github.com/scheduler/pkg/task"
)

// Store defines the persistence interface used by the scheduler and API.
type Store interface {
	SaveTask(ctx context.Context, t *task.Task) error
	GetTask(ctx context.Context, taskID string) (*task.Task, error)
	ListTasks(ctx context.Context) ([]*task.Task, error)
	DeleteTask(ctx context.Context, taskID string) error
	SaveTaskLog(ctx context.Context, log *task.TaskLog) error
	GetTaskLogs(ctx context.Context, taskID string, limit int64) ([]*task.TaskLog, error)
	Close(ctx context.Context) error
}

// Compile-time assertion that MongoDBStore satisfies Store.
var _ Store = (*MongoDBStore)(nil)

// Timeout helper exported for callers that need the default timeout.
var _ = time.Second
