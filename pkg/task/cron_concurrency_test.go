//go:build race
// +build race

package task

import (
	"fmt"
	"sync"
	"testing"
)

func TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree(t *testing.T) {
	scheduler := NewCronScheduler()
	defer scheduler.Stop()

	const workers = 8
	const iterations = 200
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(workers)
	for worker := 0; worker < workers; worker++ {
		worker := worker
		go func() {
			defer wg.Done()
			<-start
			for iteration := 0; iteration < iterations; iteration++ {
				job := &Task{
					ID:       fmt.Sprintf("worker-%d-job-%d", worker, iteration),
					CronExpr: "* * * * * *",
				}
				if err := scheduler.ScheduleTask(job, func() {}); err != nil {
					t.Errorf("schedule job: %v", err)
				}
				scheduler.UnscheduleTask(job.ID)
			}
		}()
	}
	close(start)
	wg.Wait()
}
