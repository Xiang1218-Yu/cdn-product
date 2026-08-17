package api

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/scheduler/internal/store"
	"github.com/scheduler/pkg/task"
)

type Handler struct {
	store  store.Store
	logger *zap.Logger
}

func NewHandler(store store.Store, logger *zap.Logger) *Handler {
	return &Handler{
		store:  store,
		logger: logger,
	}
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	t := task.NewTask(req.Name, req.Command, req.CronExpr)
	t.Timeout = time.Duration(req.Timeout) * time.Second
	t.MaxRetries = req.MaxRetries
	t.DependsOn = req.DependsOn

	for k, v := range req.Params {
		t.SetParam(k, v)
	}

	if err := h.store.SaveTask(r.Context(), t); err != nil {
		h.logger.Error("failed to save task", zap.Error(err))
		h.writeError(w, http.StatusInternalServerError, "failed to create task")
		return
	}

	h.writeJSON(w, http.StatusCreated, t)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		h.writeError(w, http.StatusBadRequest, "task id is required")
		return
	}

	t, err := h.store.GetTask(r.Context(), taskID)
	if err != nil {
		h.logger.Error("failed to get task", zap.Error(err))
		h.writeError(w, http.StatusNotFound, "task not found")
		return
	}

	h.writeJSON(w, http.StatusOK, t)
}

func (h *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.ListTasks(r.Context())
	if err != nil {
		h.logger.Error("failed to list tasks", zap.Error(err))
		h.writeError(w, http.StatusInternalServerError, "failed to list tasks")
		return
	}

	h.writeJSON(w, http.StatusOK, tasks)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		h.writeError(w, http.StatusBadRequest, "task id is required")
		return
	}

	if err := h.store.DeleteTask(r.Context(), taskID); err != nil {
		h.logger.Error("failed to delete task", zap.Error(err))
		h.writeError(w, http.StatusInternalServerError, "failed to delete task")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]string{"message": "task deleted"})
}

func (h *Handler) GetTaskLogs(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("task_id")
	if taskID == "" {
		h.writeError(w, http.StatusBadRequest, "task_id is required")
		return
	}

	logs, err := h.store.GetTaskLogs(r.Context(), taskID, 100)
	if err != nil {
		h.logger.Error("failed to get task logs", zap.Error(err))
		h.writeError(w, http.StatusInternalServerError, "failed to get task logs")
		return
	}

	h.writeJSON(w, http.StatusOK, logs)
}

func (h *Handler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) writeError(w http.ResponseWriter, statusCode int, message string) {
	h.writeJSON(w, statusCode, map[string]string{"error": message})
}

type CreateTaskRequest struct {
	Name       string            `json:"name"`
	Command    string            `json:"command"`
	CronExpr   string            `json:"cron_expr"`
	Params     map[string]string `json:"params"`
	Timeout    int               `json:"timeout"`
	MaxRetries int               `json:"max_retries"`
	DependsOn  []string          `json:"depends_on"`
}
