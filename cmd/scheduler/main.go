package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/scheduler/internal/config"
	"github.com/scheduler/internal/discovery"
	"github.com/scheduler/internal/lock"
	"github.com/scheduler/internal/store"
	"github.com/scheduler/pkg/api"
	"github.com/scheduler/pkg/monitor"
	"github.com/scheduler/pkg/scheduler"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "path to config file")
	flag.Parse()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	mongoStore, err := store.NewMongoDBStore(&cfg.MongoDB, logger)
	if err != nil {
		logger.Fatal("failed to connect to mongodb", zap.Error(err))
	}
	defer mongoStore.Close(context.Background())

	etcdLock, err := lock.NewEtcdLock(cfg.Etcd.Endpoints, cfg.Etcd.DialTimeout, logger)
	if err != nil {
		logger.Fatal("failed to connect to etcd", zap.Error(err))
	}
	defer etcdLock.Close()

	serviceDiscovery, err := discovery.NewServiceDiscovery(cfg.Etcd.Endpoints, cfg.Etcd.DialTimeout, logger)
	if err != nil {
		logger.Fatal("failed to create service discovery", zap.Error(err))
	}
	defer serviceDiscovery.Close()

	schedulerEngine := scheduler.NewSchedulerEngine(cfg, mongoStore, etcdLock, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := schedulerEngine.Start(ctx); err != nil {
		logger.Fatal("failed to start scheduler", zap.Error(err))
	}

	apiServer := api.NewServer(&cfg.Server.Scheduler, mongoStore, logger)

	go func() {
		if err := apiServer.Start(); err != nil {
			logger.Error("API server error", zap.Error(err))
		}
	}()

	if cfg.Prometheus.Enabled {
		go func() {
			if err := monitor.StartMetricsServer(cfg.Prometheus.Port, logger); err != nil {
				logger.Error("metrics server error", zap.Error(err))
			}
		}()
	}

	logger.Info("scheduler started successfully")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info("shutting down scheduler...")
	cancel()
	schedulerEngine.Stop()
	apiServer.Stop(context.Background())

	logger.Info("scheduler stopped")
}
