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
	reservations map[string]string // taskID -> executorID, slots held open by AcquireExecutor
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
		// The reported load may lag behind locally reserved slots, since a
		// task acquired on this scheduler has not yet surfaced in the next
		// heartbeat. Never let the bookkeeping drop below the slots this pool
		// knows are still in flight, or a capacity=1 executor could be handed a
		// second task before the first reservation is released.
		held := 0
		for _, eid := range ep.reservations {
			if eid == executorID {
				held++
			}
		}
		if load < held {
			load = held
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
			// Drop reservations pointing at the removed executor; ReleaseExecutor
			// is still safe to call later (it no-ops for unknown taskIDs), but
			// clearing them avoids a slow leak if tasks never call it back.
			for taskID, eid := range ep.reservations {
				if eid == id {
					delete(ep.reservations, taskID)
				}
			}
		}
	}
}

// AcquireExecutor atomically reserves a slot on the least-loaded executor that
// still has spare capacity. The reservation is tracked by taskID and must be
// released with ReleaseExecutor once the execution attempt finishes; until then
// the occupied slot stays counted in executor.Load so the same executor is not
// handed a second task while one is still in flight.
func (ep *ExecutorPool) AcquireExecutor(taskID string) *Executor {
	ep.mu.Lock()
	defer ep.mu.Unlock()

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

// ReleaseExecutor frees the slot reserved by the matching AcquireExecutor call.
func (ep *ExecutorPool) ReleaseExecutor(taskID string) {
	ep.mu.Lock()
	defer ep.mu.Unlock()

	executorID, ok := ep.reservations[taskID]
	if !ok {
		return
	}
	delete(ep.reservations, taskID)

	if executor, exists := ep.executors[executorID]; exists && executor.Load > 0 {
		executor.Load--
	}
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
