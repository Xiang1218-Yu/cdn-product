package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/scheduler/internal/config"
	"github.com/scheduler/internal/discovery"
	"github.com/scheduler/pkg/executor"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "path to config file")
	executorID := flag.String("id", "", "executor id")
	capacity := flag.Int("capacity", 10, "executor capacity")
	flag.Parse()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	if *executorID == "" {
		*executorID = fmt.Sprintf("executor-%d", os.Getpid())
	}

	serviceDiscovery, err := discovery.NewServiceDiscovery(cfg.Etcd.Endpoints, cfg.Etcd.DialTimeout, logger)
	if err != nil {
		logger.Fatal("failed to create service discovery", zap.Error(err))
	}
	defer serviceDiscovery.Close()

	executorServer := executor.NewExecutorServer(
		*executorID,
		&cfg.Server.Scheduler,
		*capacity,
		serviceDiscovery,
		logger,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Info("starting executor",
		zap.String("id", *executorID),
		zap.Int("capacity", *capacity),
	)

	if err := executorServer.Start(ctx); err != nil {
		logger.Fatal("failed to start executor", zap.Error(err))
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info("shutting down executor...")
	cancel()

	logger.Info("executor stopped")
}
