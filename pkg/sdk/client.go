package sdk

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/scheduler/api/proto"
)

type ExecutorClient struct {
	conn   *grpc.ClientConn
	client proto.ExecutorServiceClient
	logger *zap.Logger
}

func NewExecutorClient(address string, logger *zap.Logger) (*ExecutorClient, error) {
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(proto.Codec())),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to executor: %w", err)
	}

	return &ExecutorClient{
		conn:   conn,
		client: proto.NewExecutorServiceClient(conn),
		logger: logger,
	}, nil
}

func (ec *ExecutorClient) ExecuteTask(
	ctx context.Context,
	taskID string,
	taskName string,
	command string,
	params map[string]string,
	timeout int64,
) (*ExecutionResult, error) {
	req := &proto.TaskRequest{
		TaskId:   taskID,
		TaskName: taskName,
		Command:  command,
		Params:   params,
		Timeout:  timeout,
	}

	resp, err := ec.client.ExecuteTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("task execution failed: %w", err)
	}

	return &ExecutionResult{
		TaskID:   resp.TaskId,
		Success:  resp.Success,
		Output:   resp.Output,
		Error:    resp.Error,
		Duration: time.Duration(resp.Duration),
	}, nil
}

func (ec *ExecutorClient) GetTaskStatus(ctx context.Context, taskID string) (*TaskStatus, error) {
	req := &proto.TaskStatusRequest{
		TaskId: taskID,
	}

	resp, err := ec.client.GetTaskStatus(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get task status: %w", err)
	}

	return &TaskStatus{
		TaskID: resp.TaskId,
		Status: resp.Status,
		Output: resp.Output,
		Error:  resp.Error,
	}, nil
}

func (ec *ExecutorClient) SendHeartbeat(ctx context.Context, executorID, address string, capacity, load int32) error {
	req := &proto.HeartbeatRequest{
		ExecutorId: executorID,
		Address:    address,
		Capacity:   capacity,
		Load:       load,
	}

	resp, err := ec.client.Heartbeat(ctx, req)
	if err != nil {
		return fmt.Errorf("heartbeat failed: %w", err)
	}

	if !resp.Acknowledged {
		return fmt.Errorf("heartbeat not acknowledged")
	}

	return nil
}

func (ec *ExecutorClient) Close() error {
	return ec.conn.Close()
}

type ExecutionResult struct {
	TaskID   string
	Success  bool
	Output   string
	Error    string
	Duration time.Duration
}

type TaskStatus struct {
	TaskID string
	Status string
	Output string
	Error  string
}

type ExecutorSDK struct {
	executorID    string
	address       string
	capacity      int
	schedulerAddr string
	client        *ExecutorClient
	logger        *zap.Logger
}

func NewExecutorSDK(executorID, address string, capacity int, schedulerAddr string, logger *zap.Logger) *ExecutorSDK {
	return &ExecutorSDK{
		executorID:    executorID,
		address:       address,
		capacity:      capacity,
		schedulerAddr: schedulerAddr,
		logger:        logger,
	}
}

func (es *ExecutorSDK) Connect() error {
	client, err := NewExecutorClient(es.schedulerAddr, es.logger)
	if err != nil {
		return err
	}

	es.client = client
	return nil
}

func (es *ExecutorSDK) StartHeartbeat(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := es.client.SendHeartbeat(ctx, es.executorID, es.address, int32(es.capacity), 0); err != nil {
				es.logger.Error("heartbeat failed", zap.Error(err))
			} else {
				es.logger.Debug("heartbeat sent successfully")
			}
		}
	}
}

func (es *ExecutorSDK) Disconnect() error {
	if es.client != nil {
		return es.client.Close()
	}
	return nil
}
