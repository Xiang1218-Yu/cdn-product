package retry

import (
	"context"
	"errors"
	"github.com/scheduler/internal/config"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestRetryHonorsCancellationAndEventuallySucceeds(t *testing.T) {
	p := NewRetryPolicy(&config.RetryConfig{MaxAttempts: 3, InitialDelay: time.Millisecond, MaxDelay: 2 * time.Millisecond}, zap.NewNop())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := p.Execute(ctx, func() error { return errors.New("retry") }); !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v", err)
	}
	attempts := 0
	if err := p.Execute(context.Background(), func() error {
		attempts++
		if attempts < 3 {
			return errors.New("retry")
		}
		return nil
	}); err != nil || attempts != 3 {
		t.Fatalf("err=%v attempts=%d", err, attempts)
	}
}
