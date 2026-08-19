package executor

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/scheduler/api/proto"
	"github.com/scheduler/internal/config"
	"github.com/scheduler/internal/discovery"
)

type ExecutorServer struct {
	proto.UnimplementedExecutorServiceServer
	ID         string
	Address    string
	Capacity   int
	Load       int
	taskRunner *TaskRunner
	discovery  *discovery.ServiceDiscovery
	logger     *zap.Logger
	mu         sync.RWMutex
}

func NewExecutorServer(
	id string,
	cfg *config.SchedulerConfig,
	capacity int,
	discovery *discovery.ServiceDiscovery,
	logger *zap.Logger,
) *ExecutorServer {
	return &ExecutorServer{
		ID:         id,
		Address:    fmt.Sprintf("%s:%d", cfg.Host, cfg.GRPCPort),
		Capacity:   capacity,
		Load:       0,
		taskRunner: NewTaskRunner(),
		discovery:  discovery,
		logger:     logger,
	}
}

func (es *ExecutorServer) ExecuteTask(ctx context.Context, req *proto.TaskRequest) (*proto.TaskResponse, error) {
	es.mu.Lock()
	if es.Load >= es.Capacity {
		es.mu.Unlock()
		return &proto.TaskResponse{
			TaskId:  req.TaskId,
			Success: false,
			Error:   "executor at full capacity",
		}, nil
	}
	es.Load++
	es.mu.Unlock()

	defer func() {
		es.mu.Lock()
		es.Load--
		es.mu.Unlock()
	}()

	es.logger.Info("executing task",
		zap.String("task_id", req.TaskId),
		zap.String("task_name", req.TaskName),
	)

	startTime := time.Now()
	output, err := es.taskRunner.Run(ctx, req.TaskId, req.Command, req.Params, req.Timeout)
	duration := time.Since(startTime)

	response := &proto.TaskResponse{
		TaskId:   req.TaskId,
		Duration: int64(duration),
	}

	if err != nil {
		response.Success = false
		// Surface the timeout sentinel in the error string so callers can
		// distinguish a deadline from a generic failure without an extra RPC.
		response.Error = err.Error()
		es.logger.Error("task execution failed",
			zap.String("task_id", req.TaskId),
			zap.String("status", es.taskRunner.GetStatus(req.TaskId)),
			zap.Error(err),
		)
	} else {
		response.Success = true
		response.Output = output
		es.logger.Info("task executed successfully",
			zap.String("task_id", req.TaskId),
			zap.Duration("duration", duration),
		)
	}

	return response, nil
}

func (es *ExecutorServer) GetTaskStatus(ctx context.Context, req *proto.TaskStatusRequest) (*proto.TaskStatusResponse, error) {
	status := es.taskRunner.GetStatus(req.TaskId)

	return &proto.TaskStatusResponse{
		TaskId: req.TaskId,
		Status: status,
	}, nil
}

func (es *ExecutorServer) Heartbeat(ctx context.Context, req *proto.HeartbeatRequest) (*proto.HeartbeatResponse, error) {
	return &proto.HeartbeatResponse{Acknowledged: true}, nil
}

func (es *ExecutorServer) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", es.Address)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	proto.RegisterExecutorServiceServer(grpcServer, es)

	if err := es.discovery.Register(ctx, "executor", es.ID, es.Address, 10); err != nil {
		return err
	}

	es.logger.Info("executor server started",
		zap.String("id", es.ID),
		zap.String("address", es.Address),
	)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			es.logger.Error("grpc server error", zap.Error(err))
		}
	}()

	go es.sendHeartbeat(ctx)

	<-ctx.Done()
	grpcServer.GracefulStop()
	return nil
}

func (es *ExecutorServer) sendHeartbeat(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			es.mu.RLock()
			load := es.Load
			es.mu.RUnlock()

			es.logger.Debug("sending heartbeat",
				zap.String("executor_id", es.ID),
				zap.Int("load", load),
			)
		}
	}
}
