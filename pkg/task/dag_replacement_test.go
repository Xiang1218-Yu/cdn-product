package task

import (
	"errors"
	"testing"
)

func TestRejectedCyclicReplacementKeepsExistingTask(t *testing.T) {
	dag := NewDAG()
	worker := NewTask("worker", "consume", "")
	worker.ID = "worker"
	if err := dag.AddTask(worker); err != nil {
		t.Fatal(err)
	}

	report := NewTask("report", "aggregate", "")
	report.ID = "report"
	report.AddDependency("worker")
	if err := dag.AddTask(report); err != nil {
		t.Fatal(err)
	}

	rejected := NewTask("worker", "consume-v2", "")
	rejected.ID = "worker"
	rejected.AddDependency("report")
	if err := dag.AddTask(rejected); !errors.Is(err, ErrCycleDependency) {
		t.Fatalf("expected the replacement to be rejected for a cycle, got %v", err)
	}

	current, err := dag.GetTask("worker")
	if err != nil {
		t.Fatalf("rejected update removed the existing worker: %v", err)
	}
	if current != worker {
		t.Fatalf("rejected update replaced the existing worker: %#v", current)
	}
	if ready := dag.GetReadyTasks(); len(ready) != 1 || ready[0].ID != "worker" {
		t.Fatalf("existing worker can no longer be scheduled: %#v", ready)
	}
}
