package retry

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/scheduler/internal/config"
)

type RetryPolicy struct {
	maxAttempts  int
	initialDelay time.Duration
	maxDelay     time.Duration
	logger       *zap.Logger
}

func NewRetryPolicy(cfg *config.RetryConfig, logger *zap.Logger) *RetryPolicy {
	return &RetryPolicy{
		maxAttempts:  cfg.MaxAttempts,
		initialDelay: cfg.InitialDelay,
		maxDelay:     cfg.MaxDelay,
		logger:       logger,
	}
}

func (rp *RetryPolicy) Execute(ctx context.Context, fn func() error) error {
	var lastErr error
	delay := rp.initialDelay

	for attempt := 1; attempt <= rp.maxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err
		rp.logger.Warn("operation failed, will retry",
			zap.Int("attempt", attempt),
			zap.Int("max_attempts", rp.maxAttempts),
			zap.Error(err),
		)

		if attempt < rp.maxAttempts {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
				delay = rp.nextDelay(delay)
			}
		}
	}

	return lastErr
}

func (rp *RetryPolicy) nextDelay(currentDelay time.Duration) time.Duration {
	next := currentDelay * 2
	if next > rp.maxDelay {
		next = rp.maxDelay
	}
	return next
}

func (rp *RetryPolicy) ExecuteWithResult(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
	var lastErr error
	delay := rp.initialDelay

	for attempt := 1; attempt <= rp.maxAttempts; attempt++ {
		result, err := fn()
		if err == nil {
			return result, nil
		}

		lastErr = err
		rp.logger.Warn("operation failed, will retry",
			zap.Int("attempt", attempt),
			zap.Int("max_attempts", rp.maxAttempts),
			zap.Error(err),
		)

		if attempt < rp.maxAttempts {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
				delay = rp.nextDelay(delay)
			}
		}
	}

	return nil, lastErr
}
