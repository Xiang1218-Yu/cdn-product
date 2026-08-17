package monitor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/scheduler/pkg/task"
)

type AlertLevel string

const (
	AlertLevelInfo    AlertLevel = "info"
	AlertLevelWarning AlertLevel = "warning"
	AlertLevelError   AlertLevel = "error"
)

type Alert struct {
	ID        string
	Level     AlertLevel
	Title     string
	Message   string
	TaskID    string
	Timestamp time.Time
}

type Alerter interface {
	Send(ctx context.Context, alert *Alert) error
}

type AlertManager struct {
	alerters []Alerter
	logger   *zap.Logger
	mu       sync.RWMutex
}

func NewAlertManager(logger *zap.Logger) *AlertManager {
	return &AlertManager{
		alerters: []Alerter{},
		logger:   logger,
	}
}

func (am *AlertManager) RegisterAlerter(alerter Alerter) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.alerters = append(am.alerters, alerter)
}

func (am *AlertManager) SendAlert(ctx context.Context, alert *Alert) error {
	am.mu.RLock()
	defer am.mu.RUnlock()

	for _, alerter := range am.alerters {
		if err := alerter.Send(ctx, alert); err != nil {
			am.logger.Error("failed to send alert",
				zap.String("alert_id", alert.ID),
				zap.Error(err),
			)
			continue
		}
		am.logger.Info("alert sent successfully",
			zap.String("alert_id", alert.ID),
			zap.String("level", string(alert.Level)),
		)
	}

	return nil
}

func (am *AlertManager) CreateTaskFailureAlert(t *task.Task, err error) *Alert {
	return &Alert{
		ID:        fmt.Sprintf("task-failure-%s-%d", t.ID, time.Now().Unix()),
		Level:     AlertLevelError,
		Title:     fmt.Sprintf("Task Failed: %s", t.Name),
		Message:   fmt.Sprintf("Task %s failed with error: %v", t.ID, err),
		TaskID:    t.ID,
		Timestamp: time.Now(),
	}
}

func (am *AlertManager) CreateExecutorDownAlert(executorID, address string) *Alert {
	return &Alert{
		ID:        fmt.Sprintf("executor-down-%s-%d", executorID, time.Now().Unix()),
		Level:     AlertLevelError,
		Title:     fmt.Sprintf("Executor Down: %s", executorID),
		Message:   fmt.Sprintf("Executor %s at %s is not responding", executorID, address),
		Timestamp: time.Now(),
	}
}

func (am *AlertManager) CreateHighLoadAlert(load float64) *Alert {
	return &Alert{
		ID:        fmt.Sprintf("high-load-%d", time.Now().Unix()),
		Level:     AlertLevelWarning,
		Title:     "High Executor Load",
		Message:   fmt.Sprintf("Executor load is high: %.2f", load),
		Timestamp: time.Now(),
	}
}

type EmailAlerter struct {
	smtpHost string
	smtpPort int
	from     string
	to       []string
	logger   *zap.Logger
}

func NewEmailAlerter(smtpHost string, smtpPort int, from string, to []string, logger *zap.Logger) *EmailAlerter {
	return &EmailAlerter{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		from:     from,
		to:       to,
		logger:   logger,
	}
}

func (ea *EmailAlerter) Send(ctx context.Context, alert *Alert) error {
	ea.logger.Info("sending email alert",
		zap.String("alert_id", alert.ID),
		zap.String("level", string(alert.Level)),
	)
	return nil
}

type SlackAlerter struct {
	webhookURL string
	logger     *zap.Logger
}

func NewSlackAlerter(webhookURL string, logger *zap.Logger) *SlackAlerter {
	return &SlackAlerter{
		webhookURL: webhookURL,
		logger:     logger,
	}
}

func (sa *SlackAlerter) Send(ctx context.Context, alert *Alert) error {
	sa.logger.Info("sending slack alert",
		zap.String("alert_id", alert.ID),
		zap.String("level", string(alert.Level)),
	)
	return nil
}
