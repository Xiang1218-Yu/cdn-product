package monitor

import (
	"errors"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Metrics struct {
	tasksTotal      prometheus.Counter
	tasksSuccess    prometheus.Counter
	tasksFailed     prometheus.Counter
	tasksDuration   prometheus.Histogram
	executorsActive prometheus.Gauge
	executorsLoad   prometheus.Gauge
	schedulerQueue  prometheus.Gauge
}

func NewMetrics() *Metrics {
	m := &Metrics{
		tasksTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "scheduler_tasks_total",
			Help: "Total number of tasks submitted",
		}),
		tasksSuccess: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "scheduler_tasks_success_total",
			Help: "Total number of successful tasks",
		}),
		tasksFailed: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "scheduler_tasks_failed_total",
			Help: "Total number of failed tasks",
		}),
		tasksDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "scheduler_tasks_duration_seconds",
			Help:    "Task execution duration in seconds",
			Buckets: prometheus.DefBuckets,
		}),
		executorsActive: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "scheduler_executors_active",
			Help: "Number of active executors",
		}),
		executorsLoad: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "scheduler_executors_load",
			Help: "Current load on executors",
		}),
		schedulerQueue: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "scheduler_queue_size",
			Help: "Number of tasks in queue",
		}),
	}

	// Register each collector individually so that a duplicate registration
	// (e.g. when NewMetrics is called again in the same process, such as
	// during tests or hot reloads) reuses the already-registered collector
	// instead of panicking via MustRegister. The returned *Metrics always
	// references collectors that back the same metric families exposed on the
	// default registry, so existing usage behaviour is unchanged.
	registerOrReuse := func(c prometheus.Collector) prometheus.Collector {
		if err := prometheus.Register(c); err != nil {
			var are prometheus.AlreadyRegisteredError
			if errors.As(err, &are) {
				return are.ExistingCollector
			}
			// Any other registration error is a real problem; fall back to
			// MustRegister so it surfaces loudly rather than silently.
			panic(err)
		}
		return c
	}

	m.tasksTotal = registerOrReuse(m.tasksTotal).(prometheus.Counter)
	m.tasksSuccess = registerOrReuse(m.tasksSuccess).(prometheus.Counter)
	m.tasksFailed = registerOrReuse(m.tasksFailed).(prometheus.Counter)
	m.tasksDuration = registerOrReuse(m.tasksDuration).(prometheus.Histogram)
	m.executorsActive = registerOrReuse(m.executorsActive).(prometheus.Gauge)
	m.executorsLoad = registerOrReuse(m.executorsLoad).(prometheus.Gauge)
	m.schedulerQueue = registerOrReuse(m.schedulerQueue).(prometheus.Gauge)

	return m
}

func (m *Metrics) IncTasksTotal() {
	m.tasksTotal.Inc()
}

func (m *Metrics) IncTasksSuccess() {
	m.tasksSuccess.Inc()
}

func (m *Metrics) IncTasksFailed() {
	m.tasksFailed.Inc()
}

func (m *Metrics) RecordDuration(duration float64) {
	m.tasksDuration.Observe(duration)
}

func (m *Metrics) SetExecutorsActive(count float64) {
	m.executorsActive.Set(count)
}

func (m *Metrics) SetExecutorsLoad(load float64) {
	m.executorsLoad.Set(load)
}

func (m *Metrics) SetQueueSize(size float64) {
	m.schedulerQueue.Set(size)
}

func StartMetricsServer(port int, logger *zap.Logger) error {
	http.Handle("/metrics", promhttp.Handler())

	logger.Info("starting metrics server", zap.Int("port", port))
	return http.ListenAndServe(":2112", nil)
}
