package scheduler

import "testing"

func TestAcquireExecutorReservesCapacityUntilTaskCompletes(t *testing.T) {
	pool := NewExecutorPool()
	pool.Register(&Executor{ID: "one-slot", Capacity: 1})

	first := pool.AcquireExecutor("import-2026-08-17")
	if first == nil || first.ID != "one-slot" {
		t.Fatalf("first task was not assigned to the available executor: %#v", first)
	}
	if second := pool.AcquireExecutor("report-2026-08-17"); second != nil {
		t.Fatalf("capacity=1 admitted a second task before the first completed: %#v", second)
	}

	pool.ReleaseExecutor("import-2026-08-17")
	if next := pool.AcquireExecutor("report-2026-08-17"); next == nil || next.ID != "one-slot" {
		t.Fatalf("released slot was not reusable: %#v", next)
	}
}
