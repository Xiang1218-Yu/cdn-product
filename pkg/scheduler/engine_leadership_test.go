package scheduler

import (
	"testing"

	"go.uber.org/zap"
)

type recordingLease struct {
	released bool
}

func (lease *recordingLease) Release() error {
	lease.released = true
	return nil
}

func TestReleasedLeaderLeaseStopsScheduling(t *testing.T) {
	engine := &SchedulerEngine{isLeader: true, logger: zap.NewNop()}
	lease := &recordingLease{}

	engine.releaseLeadership(lease)

	if !lease.released {
		t.Fatal("leader lease was not released")
	}
	engine.mu.RLock()
	stillLeader := engine.isLeader
	engine.mu.RUnlock()
	if stillLeader {
		t.Fatal("released leader lease still leaves scheduler eligible to dispatch tasks")
	}
}
