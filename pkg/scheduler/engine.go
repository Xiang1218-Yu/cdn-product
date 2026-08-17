package scheduler

import (
	"context"
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/scheduler/internal/config"
	"github.com/scheduler/internal/lock"
	"github.com/scheduler/internal/store"
	"github.com/scheduler/pkg/task"
)

type SchedulerEngine struct {
	cfg           *config.Config
	store         store.Store
	lock          lock.Lock
	dag           *task.DAG
	cronScheduler *task.CronScheduler
	executorPool  *ExecutorPool
	logger        *zap.Logger
	mu            sync.RWMutex
	isLeader      bool
	nodeID        string
}

func NewSchedulerEngine(
	cfg *config.Config,
	store store.Store,
	lock lock.Lock,
	logger *zap.Logger,
) *SchedulerEngine {
	return &SchedulerEngine{
		cfg:           cfg,
		store:         store,
		lock:          lock,
		dag:           task.NewDAG(),
		cronScheduler: task.NewCronScheduler(),
		executorPool:  NewExecutorPool(),
		logger:        logger,
		nodeID:        generateNodeID(),
	}
}

func (se *SchedulerEngine) RegisterTask(t *task.Task) error {
	if t.CronExpr == "" {
		return nil
	}
	if err := se.dag.AddTask(t); err != nil {
		return err
	}
	return se.cronScheduler.ScheduleTask(t, func() {
		go se.executeTask(context.Background(), t)
	})
}

func (se *SchedulerEngine) UnregisterTask(taskID string) error {
	se.cronScheduler.UnscheduleTask(taskID)
	return se.dag.RemoveTask(taskID)
}

func (se *SchedulerEngine) Start(ctx context.Context) error {
	se.cronScheduler.Start()

	go se.electLeader(ctx)
	go se.scheduleTasks(ctx)
	go se.monitorExecutors(ctx)

	return nil
}

func (se *SchedulerEngine) electLeader(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			mutex, err := se.lock.TryLock(ctx, "/scheduler/leader", 10)
			if err != nil {
				se.mu.Lock()
				se.isLeader = false
				se.mu.Unlock()
				continue
			}

			se.mu.Lock()
			se.isLeader = true
			se.mu.Unlock()

			go func() {
				select {
				case <-ctx.Done():
					se.lock.ReleaseLock(mutex)
				case <-time.After(8 * time.Second):
					se.lock.ReleaseLock(mutex)
				}
			}()
		}
	}
}

func (se *SchedulerEngine) scheduleTasks(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			se.mu.RLock()
			isLeader := se.isLeader
			se.mu.RUnlock()

			if !isLeader {
				continue
			}

			readyTasks := se.dag.GetReadyTasks()
			for _, t := range readyTasks {
				if se.shouldHandleTask(t.ID) {
					go se.executeTask(ctx, t)
				}
			}
		}
	}
}

func (se *SchedulerEngine) shouldHandleTask(taskID string) bool {
	h := fnv.New32a()
	h.Write([]byte(taskID))
	hash := h.Sum32()

	totalNodes := se.getTotalNodes()
	nodeIndex := se.getNodeIndex()

	return int(hash)%totalNodes == nodeIndex
}

func (se *SchedulerEngine) getTotalNodes() int {
	return 3
}

func (se *SchedulerEngine) getNodeIndex() int {
	return 0
}

func (se *SchedulerEngine) executeTask(ctx context.Context, t *task.Task) {
	se.dag.UpdateTaskStatus(t.ID, task.StatusRunning)

	executor := se.executorPool.SelectExecutor(t.ID)
	if executor == nil {
		se.logger.Error("no available executor", zap.String("task_id", t.ID))
		se.dag.UpdateTaskStatus(t.ID, task.StatusFailed)
		return
	}

	result, err := se.executorPool.Execute(ctx, t)
	if err != nil {
		se.logger.Error("task execution failed",
			zap.String("task_id", t.ID),
			zap.Error(err),
		)
		se.dag.UpdateTaskStatus(t.ID, task.StatusFailed)
		se.handleFailure(ctx, t, err)
		return
	}

	if result.Success {
		se.dag.UpdateTaskStatus(t.ID, task.StatusSuccess)
		se.logger.Info("task executed successfully",
			zap.String("task_id", t.ID),
			zap.Duration("duration", result.Duration),
		)
	} else {
		se.dag.UpdateTaskStatus(t.ID, task.StatusFailed)
		se.handleFailure(ctx, t, fmt.Errorf(result.Error))
	}
}

func (se *SchedulerEngine) handleFailure(ctx context.Context, t *task.Task, err error) {
	se.logger.Error("handling task failure",
		zap.String("task_id", t.ID),
		zap.Error(err),
	)
}

func (se *SchedulerEngine) monitorExecutors(ctx context.Context) {
	ticker := time.NewTicker(se.cfg.Executor.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			se.executorPool.CleanupDeadExecutors(se.cfg.Executor.HeartbeatTimeout)
		}
	}
}

func (se *SchedulerEngine) Stop() {
	se.cronScheduler.Stop()
}

func generateNodeID() string {
	return fmt.Sprintf("scheduler-%d", time.Now().UnixNano())
}
