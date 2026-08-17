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
	executors    map[string]*Executor
	reservations map[string]string
	mu           sync.RWMutex
	logger       *zap.Logger
}

func NewExecutorPool() *ExecutorPool {
	return &ExecutorPool{
		executors:    make(map[string]*Executor),
		reservations: make(map[string]string),
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
	for taskID, reservedID := range ep.reservations {
		if reservedID == executorID {
			delete(ep.reservations, taskID)
		}
	}
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
		reserved := 0
		for _, reservedID := range ep.reservations {
			if reservedID == executorID {
				reserved++
			}
		}
		if load < reserved {
			load = reserved
		}
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
			for taskID, reservedID := range ep.reservations {
				if reservedID == id {
					delete(ep.reservations, taskID)
				}
			}
		}
	}
}

// AcquireExecutor atomically selects and reserves a slot. Repeated calls for the
// same task return its existing reservation without consuming another slot.
func (ep *ExecutorPool) AcquireExecutor(taskID string) *Executor {
	ep.mu.Lock()
	defer ep.mu.Unlock()

	if executorID, exists := ep.reservations[taskID]; exists {
		return ep.executors[executorID]
	}

	var selected *Executor
	minLoad := int(^uint(0) >> 1)
	for _, executor := range ep.executors {
		if executor.Load < executor.Capacity && executor.Load < minLoad {
			selected = executor
			minLoad = executor.Load
		}
	}
	if selected == nil {
		return nil
	}

	selected.Load++
	ep.reservations[taskID] = selected.ID
	return selected
}

// ReleaseExecutor makes an acquired slot available after completion or failover.
func (ep *ExecutorPool) ReleaseExecutor(taskID string) {
	ep.mu.Lock()
	defer ep.mu.Unlock()

	executorID, exists := ep.reservations[taskID]
	if !exists {
		return
	}
	if executor, exists := ep.executors[executorID]; exists && executor.Load > 0 {
		executor.Load--
	}
	delete(ep.reservations, taskID)
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
