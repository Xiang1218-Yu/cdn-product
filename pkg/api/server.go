package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/scheduler/internal/config"
	"github.com/scheduler/internal/store"
)

type Server struct {
	httpServer *http.Server
	handler    *Handler
	logger     *zap.Logger
}

func NewServer(cfg *config.SchedulerConfig, store store.Store, logger *zap.Logger) *Server {
	handler := NewHandler(store, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/tasks", handler.CreateTask)
	mux.HandleFunc("/api/v1/tasks/list", handler.ListTasks)
	mux.HandleFunc("/api/v1/tasks/get", handler.GetTask)
	mux.HandleFunc("/api/v1/tasks/delete", handler.DeleteTask)
	mux.HandleFunc("/api/v1/logs", handler.GetTaskLogs)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		handler:    handler,
		logger:     logger,
	}
}

func (s *Server) Start() error {
	s.logger.Info("starting API server", zap.String("address", s.httpServer.Addr))
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("stopping API server")
	return s.httpServer.Shutdown(ctx)
}
