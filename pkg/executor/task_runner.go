package executor

import (
	"context"
	"errors"
	"sync"
	"time"
)

// taskTimeout is the floor for a per-task timeout. A request with a zero or
// negative Timeout runs under the caller's deadline only; a sub-floor timeout
// would otherwise produce an already-expired context and a spurious failure.
const taskTimeout = 1 * time.Millisecond

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

// Status constants recorded on RunningTask. Distinct from pkg/task.TaskStatus
// because the executor reports an out-of-band "timed_out" terminal state so
// the scheduler can tell a task that blew its own deadline apart from one
// that simply returned an error.
const (
	StatusRunning  = "running"
	StatusSuccess  = "success"
	StatusFailed   = "failed"
	StatusTimedOut = "timed_out"
)

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

// Run executes the user command under a per-task deadline derived from
// timeout (in nanoseconds). The deadline is applied here, at the execution
// boundary, rather than relied upon to be carried by the caller's context —
// otherwise a short task timeout only takes effect when the (typically much
// longer) request context expires. The derived context is canceled as soon
// as Run returns so the runner function cannot outlive the request.
func (tr *TaskRunner) Run(ctx context.Context, taskID, command string, params map[string]string, timeout int64) (string, error) {
	tr.mu.Lock()
	tr.tasks[taskID] = &RunningTask{ID: taskID, Command: command, Status: StatusRunning, StartTime: time.Now()}
	tr.mu.Unlock()

	runCtx, parent, cancel := taskContext(ctx, timeout)
	defer cancel()

	output, err := tr.run(runCtx, command, params)

	tr.mu.Lock()
	defer tr.mu.Unlock()
	running := tr.tasks[taskID]
	if err != nil {
		// Distinguish a task that blew its own deadline from one that simply
		// returned an error. DeadlineExceeded can come from the per-task
		// context or from an upstream deadline/cancellation; the task timeout
		// is the culprit only when the task context expired while the parent
		// had not yet been canceled or deadlined.
		running.Status = StatusFailed
		if isTaskTimeout(runCtx, parent, err) {
			running.Status = StatusTimedOut
		}
		running.Error = err.Error()
		return "", err
	}
	running.Status = StatusSuccess
	running.Output = output
	return output, nil
}

// taskContext derives a child context whose deadline is the earlier of the
// caller's deadline and the per-task timeout. It returns the parent context
// unchanged so the caller can compare the two to attribute a deadline to the
// right source. A non-positive timeout means "no per-task deadline": the
// caller's context governs and the child is a plain cancel wrapper.
func taskContext(ctx context.Context, timeout int64) (context.Context, context.Context, context.CancelFunc) {
	if timeout <= 0 {
		ctx, cancel := context.WithCancel(ctx)
		return ctx, ctx, cancel
	}
	d := time.Duration(timeout)
	if d < taskTimeout {
		d = taskTimeout
	}
	runCtx, cancel := context.WithTimeout(ctx, d)
	return runCtx, ctx, cancel
}

// isTaskTimeout reports whether err is a deadline-exceeded that came from the
// per-task context rather than an upstream cancellation or parent deadline.
func isTaskTimeout(taskCtx, parent context.Context, err error) bool {
	if !errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	// The task context must be done by its own deadline. If the parent was
	// already canceled/deadlined first, the cause is upstream.
	if parent.Err() != nil {
		return false
	}
	if dl, ok := taskCtx.Deadline(); ok {
		return !time.Now().Before(dl)
	}
	return true
}

func (tr *TaskRunner) GetStatus(taskID string) string {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	if task, exists := tr.tasks[taskID]; exists {
		return task.Status
	}
	return "unknown"
}
