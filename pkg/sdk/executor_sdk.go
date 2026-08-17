package sdk

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type TaskHandler func(ctx context.Context, command string, params map[string]string) (string, error)

type SimpleExecutor struct {
	id       string
	address  string
	capacity int
	handler  TaskHandler
	client   *ExecutorClient
	logger   *zap.Logger
}

func NewSimpleExecutor(id, address string, capacity int, handler TaskHandler, logger *zap.Logger) *SimpleExecutor {
	return &SimpleExecutor{
		id:       id,
		address:  address,
		capacity: capacity,
		handler:  handler,
		logger:   logger,
	}
}

func (se *SimpleExecutor) Execute(ctx context.Context, taskID, taskName, command string, params map[string]string, timeout int64) (*ExecutionResult, error) {
	se.logger.Info("executing task",
		zap.String("task_id", taskID),
		zap.String("task_name", taskName),
	)

	output, err := se.handler(ctx, command, params)
	if err != nil {
		return &ExecutionResult{
			TaskID:  taskID,
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &ExecutionResult{
		TaskID:  taskID,
		Success: true,
		Output:  output,
	}, nil
}

func (se *SimpleExecutor) Register(schedulerAddr string) error {
	client, err := NewExecutorClient(schedulerAddr, se.logger)
	if err != nil {
		return fmt.Errorf("failed to connect to scheduler: %w", err)
	}

	se.client = client
	return nil
}

func ExampleUsage() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	handler := func(ctx context.Context, command string, params map[string]string) (string, error) {
		return fmt.Sprintf("Executed: %s", command), nil
	}

	executor := NewSimpleExecutor("executor-1", "localhost:9090", 10, handler, logger)

	if err := executor.Register("localhost:9090"); err != nil {
		logger.Fatal("failed to register executor", zap.Error(err))
	}

	logger.Info("executor registered successfully")
}
