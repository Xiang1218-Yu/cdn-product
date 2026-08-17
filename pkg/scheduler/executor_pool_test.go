package scheduler

import (
	"testing"
	"time"
)

func TestExecutorPoolPrefersAvailableLeastLoadedExecutor(t *testing.T) {
	p := NewExecutorPool()
	p.Register(&Executor{ID: "busy", Capacity: 2, Load: 2})
	p.Register(&Executor{ID: "best", Capacity: 3, Load: 1})
	p.Register(&Executor{ID: "other", Capacity: 3, Load: 2})
	if got := p.SelectExecutor("task"); got == nil || got.ID != "best" {
		t.Fatalf("selected %#v", got)
	}
	p.UpdateHeartbeat("best", 3)
	if got := p.SelectExecutor("task"); got == nil || got.ID != "other" {
		t.Fatalf("selected %#v", got)
	}
	p.CleanupDeadExecutors(-time.Second)
	if got := p.SelectExecutor("task"); got != nil {
		t.Fatalf("stale executor survived: %#v", got)
	}
}
