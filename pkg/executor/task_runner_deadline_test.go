package executor

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/scheduler/api/proto"
	"github.com/scheduler/internal/config"
	"go.uber.org/zap"
)

func TestExecutorHonorsRequestDeadlineAndRecordsTimeout(t *testing.T) {
	server := NewExecutorServer("executor-timeout", &config.SchedulerConfig{Host: "127.0.0.1", GRPCPort: 9000}, 1, nil, zap.NewNop())
	server.taskRunner.run = func(ctx context.Context, command string, params map[string]string) (string, error) {
		<-ctx.Done()
		return "", ctx.Err()
	}

	parent, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	started := time.Now()
	response, err := server.ExecuteTask(parent, &proto.TaskRequest{
		TaskId:   "deadline-task",
		TaskName: "deadline-task",
		Command:  "sleep",
		Timeout:  int64(8 * time.Millisecond),
	})
	elapsed := time.Since(started)

	if err != nil {
		t.Fatalf("ExecuteTask returned transport error: %v", err)
	}
	if response.Success || !strings.Contains(response.Error, context.DeadlineExceeded.Error()) {
		t.Fatalf("response=%#v, want a deadline failure", response)
	}
	if elapsed > 45*time.Millisecond {
		t.Fatalf("request timeout was ignored: execution lasted %s", elapsed)
	}
	if status := server.taskRunner.GetStatus("deadline-task"); status != "timed_out" {
		t.Fatalf("status=%q, want timed_out", status)
	}
	if server.Load != 0 {
		t.Fatalf("executor load leaked after timeout: %d", server.Load)
	}
}
