package task

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestConcurrentDispatchClaimsPendingTaskOnlyOnce(t *testing.T) {
	dag := NewDAG()
	pending := NewTask("refresh", "refresh-cache", "")
	pending.ID = "refresh"
	if err := dag.AddTask(pending); err != nil {
		t.Fatal(err)
	}

	start := make(chan struct{})
	var accepted atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			if _, ok := dag.ClaimReadyTask("refresh"); ok {
				accepted.Add(1)
			}
		}()
	}
	close(start)
	wg.Wait()

	if got := accepted.Load(); got != 1 {
		t.Fatalf("pending task was claimed %d times; concurrent scheduler ticks must dispatch it once", got)
	}

	got, err := dag.GetTask("refresh")
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != StatusRunning {
		t.Fatalf("claimed task status = %q, want %q", got.Status, StatusRunning)
	}
}
