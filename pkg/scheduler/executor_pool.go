package scheduler

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/scheduler/pkg/task"
)

type Executor struct {
	ID       string
	Address  string
	Capacity int
	Load     int
	LastSeen time.Time
}

type ExecutorPool struct {
	executors map[string]*Executor
	mu        sync.RWMutex
	logger    *zap.Logger
}

func NewExecutorPool() *ExecutorPool {
	return &ExecutorPool{
		executors: make(map[string]*Executor),
	}
}

func (ep *ExecutorPool) Register(executor *Executor) {
	ep.mu.Lock()
	defer ep.mu.Unlock()

	executor.LastSeen = time.Now()
	ep.executors[executor.ID] = executor
}

func (ep *ExecutorPool) Deregister(executorID string) {
	ep.mu.Lock()
	defer ep.mu.Unlock()

	delete(ep.executors, executorID)
}

func (ep *ExecutorPool) SelectExecutor(taskID string) *Executor {
	ep.mu.RLock()
	defer ep.mu.RUnlock()

	var selected *Executor
	minLoad := int(^uint(0) >> 1)

	for _, executor := range ep.executors {
		if executor.Load < executor.Capacity && executor.Load < minLoad {
			selected = executor
			minLoad = executor.Load
		}
	}

	return selected
}

func (ep *ExecutorPool) UpdateHeartbeat(executorID string, load int) {
	ep.mu.Lock()
	defer ep.mu.Unlock()

	if executor, exists := ep.executors[executorID]; exists {
		executor.Load = load
		executor.LastSeen = time.Now()
	}
}

func (ep *ExecutorPool) CleanupDeadExecutors(timeout time.Duration) {
	ep.mu.Lock()
	defer ep.mu.Unlock()

	now := time.Now()
	for id, executor := range ep.executors {
		if now.Sub(executor.LastSeen) > timeout {
			delete(ep.executors, id)
		}
	}
}

// AcquireExecutor chooses a candidate for a task. The scheduler calls ReleaseExecutor
// after the execution attempt has finished.
func (ep *ExecutorPool) AcquireExecutor(taskID string) *Executor {
	return ep.SelectExecutor(taskID)
}

func (ep *ExecutorPool) ReleaseExecutor(taskID string) {
}

func (ep *ExecutorPool) Execute(ctx context.Context, t *task.Task) (*ExecutionResult, error) {
	return nil, nil
}

type ExecutionResult struct {
	Success  bool
	Output   string
	Error    string
	Duration time.Duration
}
