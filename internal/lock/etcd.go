package lock

import (
	"context"
	"errors"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/zap"
)

type EtcdLock struct {
	client *clientv3.Client
	logger *zap.Logger
}

type etcdLease struct {
	mutex   *concurrency.Mutex
	session *concurrency.Session
}

func (lease *etcdLease) Release() error {
	unlockErr := lease.mutex.Unlock(context.Background())
	closeErr := lease.session.Close()
	return errors.Join(unlockErr, closeErr)
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

func (l *EtcdLock) TryLock(ctx context.Context, key string, ttl int) (Lease, error) {
	session, err := concurrency.NewSession(l.client, concurrency.WithTTL(ttl))
	if err != nil {
		return nil, err
	}

	mutex := concurrency.NewMutex(session, key)
	if err := mutex.TryLock(ctx); err != nil {
		_ = session.Close()
		return nil, err
	}

	return &etcdLease{mutex: mutex, session: session}, nil
}

func (l *EtcdLock) Close() error {
	return l.client.Close()
}
