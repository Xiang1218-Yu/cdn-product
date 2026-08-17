package task

import (
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type CronScheduler struct {
	parser cron.Parser
	jobs   map[string]cron.EntryID
	cron   *cron.Cron
	mu     sync.Mutex
}

func NewCronScheduler() *CronScheduler {
	return &CronScheduler{
		parser: cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow),
		jobs:   make(map[string]cron.EntryID),
		cron:   cron.New(cron.WithParser(cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow))),
	}
}

func (cs *CronScheduler) Start() {
	cs.cron.Start()
}

func (cs *CronScheduler) Stop() {
	cs.cron.Stop()
}

func (cs *CronScheduler) ScheduleTask(task *Task, runFunc func()) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	entryID, err := cs.cron.AddFunc(task.CronExpr, runFunc)
	if err != nil {
		return err
	}

	if previous, exists := cs.jobs[task.ID]; exists {
		cs.cron.Remove(previous)
	}
	cs.jobs[task.ID] = entryID
	return nil
}

func (cs *CronScheduler) UnscheduleTask(taskID string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if entryID, exists := cs.jobs[taskID]; exists {
		cs.cron.Remove(entryID)
		delete(cs.jobs, taskID)
	}
}

func (cs *CronScheduler) GetNextRunTime(cronExpr string) (time.Time, error) {
	schedule, err := cs.parser.Parse(cronExpr)
	if err != nil {
		return time.Time{}, err
	}

	return schedule.Next(time.Now()), nil
}
