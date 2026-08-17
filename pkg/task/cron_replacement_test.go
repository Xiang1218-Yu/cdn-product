package task

import "testing"

func TestReschedulingSameTaskReplacesPreviousCronEntry(t *testing.T) {
	scheduler := NewCronScheduler()
	defer scheduler.Stop()

	job := NewTask("daily-report", "report", "*/5 * * * * *")
	job.ID = "daily-report"
	if err := scheduler.ScheduleTask(job, func() {}); err != nil {
		t.Fatal(err)
	}

	job.CronExpr = "*/15 * * * * *"
	if err := scheduler.ScheduleTask(job, func() {}); err != nil {
		t.Fatal(err)
	}

	if entries := scheduler.cron.Entries(); len(entries) != 1 {
		t.Fatalf("rescheduling one task left %d active cron entries, want 1", len(entries))
	}
}
