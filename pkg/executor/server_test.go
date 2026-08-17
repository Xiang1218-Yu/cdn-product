package executor

import (
	"context"
	"github.com/scheduler/api/proto"
	"github.com/scheduler/internal/config"
	"go.uber.org/zap"
	"testing"
)

func TestExecutorRejectsRequestsAtCapacity(t *testing.T) {
	s := NewExecutorServer("e1", &config.SchedulerConfig{Host: "127.0.0.1", GRPCPort: 9000}, 1, nil, zap.NewNop())
	s.Load = 1
	r, err := s.ExecuteTask(context.Background(), &proto.TaskRequest{TaskId: "t1", TaskName: "daily"})
	if err != nil || r.Success || r.Error != "executor at full capacity" {
		t.Fatalf("response=%#v err=%v", r, err)
	}
}
