package scheduler

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/scheduler/pkg/task"
)

type FailoverManager struct {
	runningTasks map[string]*TaskExecution
	mu           sync.RWMutex
	logger       *zap.Logger
}

type TaskExecution struct {
	TaskID     string
	ExecutorID string
	StartTime  time.Time
	LastUpdate time.Time
}

func NewFailoverManager(logger *zap.Logger) *FailoverManager {
	return &FailoverManager{
		runningTasks: make(map[string]*TaskExecution),
		logger:       logger,
	}
}

func (fm *FailoverManager) RecordExecution(taskID, executorID string) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fm.runningTasks[taskID] = &TaskExecution{
		TaskID:     taskID,
		ExecutorID: executorID,
		StartTime:  time.Now(),
		LastUpdate: time.Now(),
	}
}

func (fm *FailoverManager) CompleteExecution(taskID string) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	delete(fm.runningTasks, taskID)
}

func (fm *FailoverManager) HandleExecutorFailure(ctx context.Context, executorID string, dag *task.DAG, executorPool *ExecutorPool) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	for taskID, execution := range fm.runningTasks {
		if execution.ExecutorID == executorID {
			fm.logger.Warn("handling executor failure for task",
				zap.String("task_id", taskID),
				zap.String("executor_id", executorID),
			)

			dag.UpdateTaskStatus(taskID, task.StatusPending)

			newExecutor := executorPool.SelectExecutor(taskID)
			if newExecutor == nil {
				fm.logger.Warn("no executor available for failed task",
					zap.String("task_id", taskID),
				)
				continue
			}

			execution.ExecutorID = newExecutor.ID
			execution.LastUpdate = time.Now()
			fm.logger.Info("reassigning task to new executor",
				zap.String("task_id", taskID),
				zap.String("new_executor_id", newExecutor.ID),
			)
		}
	}
}

func (fm *FailoverManager) GetRunningTasks() map[string]*TaskExecution {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	result := make(map[string]*TaskExecution)
	for k, v := range fm.runningTasks {
		result[k] = v
	}
	return result
}
