package log

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/scheduler/pkg/task"
)

type TaskLogger struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

func NewTaskLogger(db *mongo.Database, logger *zap.Logger) *TaskLogger {
	return &TaskLogger{
		collection: db.Collection("task_logs"),
		logger:     logger,
	}
}

func (tl *TaskLogger) LogExecution(ctx context.Context, taskLog *task.TaskLog) error {
	taskLog.ID = generateLogID()

	_, err := tl.collection.InsertOne(ctx, taskLog)
	if err != nil {
		tl.logger.Error("failed to log task execution",
			zap.String("task_id", taskLog.TaskID),
			zap.Error(err),
		)
		return err
	}

	tl.logger.Info("task execution logged",
		zap.String("task_id", taskLog.TaskID),
		zap.String("status", taskLog.Status),
	)

	return nil
}

func (tl *TaskLogger) GetLogs(ctx context.Context, taskID string, limit int64) ([]*task.TaskLog, error) {
	opts := bson.M{"task_id": taskID}

	cursor, err := tl.collection.Find(ctx, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*task.TaskLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return logs, nil
}

func (tl *TaskLogger) GetRecentLogs(ctx context.Context, limit int64) ([]*task.TaskLog, error) {
	cursor, err := tl.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*task.TaskLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return logs, nil
}

func generateLogID() string {
	return time.Now().Format("20060102150405")
}
