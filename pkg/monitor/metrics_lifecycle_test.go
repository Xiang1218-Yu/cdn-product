package monitor

import "testing"

func TestNewMetricsCanBeCreatedMoreThanOnce(t *testing.T) {
	first := NewMetrics()
	second := NewMetrics()
	first.IncTasksTotal()
	second.IncTasksSuccess()
}
