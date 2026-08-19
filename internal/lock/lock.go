package lock

import "context"

// Lease represents one acquired distributed leadership lease.
type Lease interface {
	Release() error
}

// Lock defines the distributed lock interface used by the scheduler.
type Lock interface {
	TryLock(ctx context.Context, key string, ttl int) (Lease, error)
	Close() error
}

// Compile-time assertion that EtcdLock satisfies Lock.
var _ Lock = (*EtcdLock)(nil)
