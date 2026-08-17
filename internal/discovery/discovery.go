package discovery

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type ServiceDiscovery struct {
	client *clientv3.Client
	logger *zap.Logger
}

func NewServiceDiscovery(endpoints []string, dialTimeout time.Duration, logger *zap.Logger) (*ServiceDiscovery, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return nil, err
	}

	return &ServiceDiscovery{
		client: cli,
		logger: logger,
	}, nil
}

func (sd *ServiceDiscovery) Register(ctx context.Context, serviceName, instanceID, address string, ttl int64) error {
	key := "/services/" + serviceName + "/" + instanceID

	lease, err := sd.client.Grant(ctx, ttl)
	if err != nil {
		return err
	}

	_, err = sd.client.Put(ctx, key, address, clientv3.WithLease(lease.ID))
	if err != nil {
		return err
	}

	keepAliveChan, err := sd.client.KeepAlive(ctx, lease.ID)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case ka := <-keepAliveChan:
				if ka == nil {
					sd.logger.Warn("keep alive channel closed")
					return
				}
			}
		}
	}()

	return nil
}

func (sd *ServiceDiscovery) Discover(ctx context.Context, serviceName string) ([]string, error) {
	resp, err := sd.client.Get(ctx, "/services/"+serviceName+"/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var addresses []string
	for _, kv := range resp.Kvs {
		addresses = append(addresses, string(kv.Value))
	}

	return addresses, nil
}

func (sd *ServiceDiscovery) Deregister(ctx context.Context, serviceName, instanceID string) error {
	key := "/services/" + serviceName + "/" + instanceID
	_, err := sd.client.Delete(ctx, key)
	return err
}

func (sd *ServiceDiscovery) Watch(ctx context.Context, serviceName string) clientv3.WatchChan {
	prefix := "/services/" + serviceName + "/"
	return sd.client.Watch(ctx, prefix, clientv3.WithPrefix())
}

func (sd *ServiceDiscovery) Close() error {
	return sd.client.Close()
}
