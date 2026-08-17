package store

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/scheduler/internal/config"
	"github.com/scheduler/pkg/task"
)

type MongoDBStore struct {
	client   *mongo.Client
	database *mongo.Database
	logger   *zap.Logger
}

func NewMongoDBStore(cfg *config.MongoDBConfig, logger *zap.Logger) (*MongoDBStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &MongoDBStore{
		client:   client,
		database: client.Database(cfg.Database),
		logger:   logger,
	}, nil
}

func (s *MongoDBStore) SaveTask(ctx context.Context, task *task.Task) error {
	collection := s.database.Collection("tasks")

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": task.ID},
		bson.M{"$set": task},
		options.Update().SetUpsert(true),
	)

	return err
}

func (s *MongoDBStore) GetTask(ctx context.Context, taskID string) (*task.Task, error) {
	collection := s.database.Collection("tasks")

	var t task.Task
	err := collection.FindOne(ctx, bson.M{"_id": taskID}).Decode(&t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (s *MongoDBStore) ListTasks(ctx context.Context) ([]*task.Task, error) {
	collection := s.database.Collection("tasks")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []*task.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *MongoDBStore) DeleteTask(ctx context.Context, taskID string) error {
	collection := s.database.Collection("tasks")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": taskID})
	return err
}

func (s *MongoDBStore) SaveTaskLog(ctx context.Context, log *task.TaskLog) error {
	collection := s.database.Collection("task_logs")

	_, err := collection.InsertOne(ctx, log)
	return err
}

func (s *MongoDBStore) GetTaskLogs(ctx context.Context, taskID string, limit int64) ([]*task.TaskLog, error) {
	collection := s.database.Collection("task_logs")

	opts := options.Find().
		SetSort(bson.D{{Key: "start_time", Value: -1}}).
		SetLimit(limit)

	cursor, err := collection.Find(ctx, bson.M{"task_id": taskID}, opts)
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

func (s *MongoDBStore) Close(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}
