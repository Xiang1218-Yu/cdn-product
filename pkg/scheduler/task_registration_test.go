package scheduler

import (
	"testing"

	"go.uber.org/zap"

	"github.com/scheduler/pkg/task"
)

func TestRegisterTaskKeepsManualTaskReadyForScheduler(t *testing.T) {
	engine := &SchedulerEngine{
		dag:           task.NewDAG(),
		cronScheduler: task.NewCronScheduler(),
		logger:        zap.NewNop(),
	}
	defer engine.cronScheduler.Stop()

	manual := task.NewTask("cache-warm", "warm-cache", "")
	manual.ID = "cache-warm"
	if err := engine.RegisterTask(manual); err != nil {
		t.Fatalf("register manual task: %v", err)
	}

	ready := engine.dag.GetReadyTasks()
	if len(ready) != 1 || ready[0].ID != manual.ID {
		t.Fatalf("manual API task was not available to the scheduler: %#v", ready)
	}
}
