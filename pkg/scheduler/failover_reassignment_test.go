package scheduler

import (
	"context"
	"testing"

	"go.uber.org/zap"

	"github.com/scheduler/pkg/task"
)

func TestFailoverRecordsReplacementExecutor(t *testing.T) {
	dag := task.NewDAG()
	report := task.NewTask("report", "aggregate", "")
	report.ID = "report"
	if err := dag.AddTask(report); err != nil {
		t.Fatal(err)
	}
	if err := dag.UpdateTaskStatus(report.ID, task.StatusRunning); err != nil {
		t.Fatal(err)
	}

	pool := NewExecutorPool()
	pool.Register(&Executor{ID: "survivor", Capacity: 1})
	failover := NewFailoverManager(zap.NewNop())
	failover.RecordExecution(report.ID, "dead")

	failover.HandleExecutorFailure(context.Background(), "dead", dag, pool)

	recorded := failover.GetRunningTasks()[report.ID]
	if recorded == nil || recorded.ExecutorID != "survivor" {
		t.Fatalf("failover selected survivor but kept stale execution record: %#v", recorded)
	}
}
