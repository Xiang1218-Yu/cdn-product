package task

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID         string            `bson:"_id" json:"id"`
	Name       string            `bson:"name" json:"name"`
	Command    string            `bson:"command" json:"command"`
	CronExpr   string            `bson:"cron_expr" json:"cron_expr"`
	Params     map[string]string `bson:"params" json:"params"`
	Timeout    time.Duration     `bson:"timeout" json:"timeout"`
	MaxRetries int               `bson:"max_retries" json:"max_retries"`
	DependsOn  []string          `bson:"depends_on" json:"depends_on"`
	Status     TaskStatus        `bson:"status" json:"status"`
	CreatedAt  time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time         `bson:"updated_at" json:"updated_at"`
}

type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusSuccess   TaskStatus = "success"
	StatusFailed    TaskStatus = "failed"
	StatusCancelled TaskStatus = "cancelled"
)

type TaskLog struct {
	ID        string    `bson:"_id" json:"id"`
	TaskID    string    `bson:"task_id" json:"task_id"`
	StartTime time.Time `bson:"start_time" json:"start_time"`
	EndTime   time.Time `bson:"end_time" json:"end_time"`
	Status    string    `bson:"status" json:"status"`
	Output    string    `bson:"output" json:"output"`
	Error     string    `bson:"error" json:"error"`
	Retry     int       `bson:"retry" json:"retry"`
}

func NewTask(name, command, cronExpr string) *Task {
	return &Task{
		ID:        uuid.New().String(),
		Name:      name,
		Command:   command,
		CronExpr:  cronExpr,
		Status:    StatusPending,
		Params:    make(map[string]string),
		DependsOn: []string{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (t *Task) AddDependency(taskID string) {
	t.DependsOn = append(t.DependsOn, taskID)
	t.UpdatedAt = time.Now()
}

func (t *Task) SetParam(key, value string) {
	t.Params[key] = value
	t.UpdatedAt = time.Now()
}
