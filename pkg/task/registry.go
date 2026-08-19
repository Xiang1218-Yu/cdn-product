package task

import (
	"sync"

	"github.com/robfig/cron/v3"
)

// cronJobRegistry serializes access to the scheduler's task-to-entry mapping.
type cronJobRegistry struct {
	mu      sync.Mutex
	entries map[string]cron.EntryID
}

func newCronJobRegistry() *cronJobRegistry {
	return &cronJobRegistry{entries: make(map[string]cron.EntryID)}
}

func (r *cronJobRegistry) add(taskID string, entryID cron.EntryID) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[taskID] = entryID
}

func (r *cronJobRegistry) take(taskID string) (cron.EntryID, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	entryID, exists := r.entries[taskID]
	if exists {
		delete(r.entries, taskID)
	}
	return entryID, exists
}
