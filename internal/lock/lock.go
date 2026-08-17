package lock

import (
	"context"

	"go.etcd.io/etcd/client/v3/concurrency"
)

// Lock defines the distributed lock interface used by the scheduler.
type Lock interface {
	AcquireLock(ctx context.Context, key string, ttl int) (*concurrency.Mutex, error)
	ReleaseLock(mutex *concurrency.Mutex) error
	TryLock(ctx context.Context, key string, ttl int) (*concurrency.Mutex, error)
	Close() error
}

// Compile-time assertion that EtcdLock satisfies Lock.
var _ Lock = (*EtcdLock)(nil)
