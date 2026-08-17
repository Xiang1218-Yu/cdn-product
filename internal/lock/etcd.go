package lock

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/zap"
)

type EtcdLock struct {
	client *clientv3.Client
	logger *zap.Logger
}

func NewEtcdLock(endpoints []string, dialTimeout time.Duration, logger *zap.Logger) (*EtcdLock, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return nil, err
	}

	return &EtcdLock{
		client: cli,
		logger: logger,
	}, nil
}

func (l *EtcdLock) AcquireLock(ctx context.Context, key string, ttl int) (*concurrency.Mutex, error) {
	session, err := concurrency.NewSession(l.client, concurrency.WithTTL(ttl))
	if err != nil {
		return nil, err
	}

	mutex := concurrency.NewMutex(session, key)
	if err := mutex.Lock(ctx); err != nil {
		return nil, err
	}

	return mutex, nil
}

func (l *EtcdLock) ReleaseLock(mutex *concurrency.Mutex) error {
	return mutex.Unlock(context.Background())
}

func (l *EtcdLock) TryLock(ctx context.Context, key string, ttl int) (*concurrency.Mutex, error) {
	session, err := concurrency.NewSession(l.client, concurrency.WithTTL(ttl))
	if err != nil {
		return nil, err
	}

	mutex := concurrency.NewMutex(session, key)
	if err := mutex.TryLock(ctx); err != nil {
		return nil, err
	}

	return mutex, nil
}

func (l *EtcdLock) Close() error {
	return l.client.Close()
}
