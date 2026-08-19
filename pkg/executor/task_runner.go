package executor

import (
	"context"
	"sync"
	"time"
)

type commandRunner func(context.Context, string, map[string]string) (string, error)

type TaskRunner struct {
	tasks map[string]*RunningTask
	mu    sync.RWMutex
	run   commandRunner
}

type RunningTask struct {
	ID        string
	Command   string
	Status    string
	Output    string
	Error     string
	StartTime time.Time
}

func NewTaskRunner() *TaskRunner {
	return &TaskRunner{
		tasks: make(map[string]*RunningTask),
		run: func(ctx context.Context, command string, params map[string]string) (string, error) {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			default:
				return "task executed successfully", nil
			}
		},
	}
}

func (tr *TaskRunner) Run(ctx context.Context, taskID, command string, params map[string]string) (string, error) {
	tr.mu.Lock()
	tr.tasks[taskID] = &RunningTask{ID: taskID, Command: command, Status: "running", StartTime: time.Now()}
	tr.mu.Unlock()

	output, err := tr.run(ctx, command, params)
	tr.mu.Lock()
	defer tr.mu.Unlock()
	running := tr.tasks[taskID]
	if err != nil {
		running.Status = "failed"
		running.Error = err.Error()
		return "", err
	}
	running.Status = "success"
	running.Output = output
	return output, nil
}

func (tr *TaskRunner) GetStatus(taskID string) string {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	if task, exists := tr.tasks[taskID]; exists {
		return task.Status
	}
	return "unknown"
}
