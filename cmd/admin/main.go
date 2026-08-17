package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/scheduler/internal/config"
	"github.com/scheduler/internal/store"
	"github.com/scheduler/pkg/api"
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

	apiServer := api.NewServer(&cfg.Server.Scheduler, mongoStore, logger)

	logger.Info("starting admin server")

	go func() {
		if err := apiServer.Start(); err != nil {
			logger.Error("API server error", zap.Error(err))
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info("shutting down admin server...")
	apiServer.Stop(context.Background())

	logger.Info("admin server stopped")
}
