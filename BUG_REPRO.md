# Bug 复现说明

## Bug 是什么
监控指标构造函数每次都会向 Prometheus 默认注册表注册同名指标。在同一进程中第二次构造监控对象时，默认注册表拒绝重复 Collector，`MustRegister` 随即触发 panic。

## 如何触发
在同一个测试进程中连续两次调用监控指标构造函数。第二次调用会尝试注册已经存在的同名指标，并导致测试或热重载流程崩溃。

## 运行指令
```bash
go test ./pkg/monitor -run '^TestNewMetricsCanBeCreatedMoreThanOnce$' -count=1
```

## 错误信息
第二次初始化触发以下 panic：
```text
panic: duplicate metrics collector registration attempted
```

## 错误堆栈
以下为上述命令在含 Bug 环境中的完整原始输出：
```text
--- FAIL: TestNewMetricsCanBeCreatedMoreThanOnce (0.00s)
panic: duplicate metrics collector registration attempted [recovered, repanicked]

goroutine 21 [running]:
testing.tRunner.func1.2({0x101359cc0, 0x3a2a6e1248c0})
	/opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:1974 +0x1a0
testing.tRunner.func1()
	/opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:1977 +0x318
panic({0x101359cc0?, 0x3a2a6e1248c0?})
	/opt/homebrew/Cellar/go/1.26.5/libexec/src/runtime/panic.go:860 +0x12c
github.com/prometheus/client_golang/prometheus.(*Registry).MustRegister(0x101431780, {0x3a2a6e1267e0?, 0x7, 0x0?})
	/Users/tog/go/pkg/mod/github.com/prometheus/client_golang@v1.17.0/prometheus/registry.go:405 +0x78
github.com/prometheus/client_golang/prometheus.MustRegister(...)
	/Users/tog/go/pkg/mod/github.com/prometheus/client_golang@v1.17.0/prometheus/registry.go:177
github.com/scheduler/pkg/monitor.NewMetrics()
	/Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__010/env/pkg/monitor/metrics.go:54 +0x554
github.com/scheduler/pkg/monitor.TestNewMetricsCanBeCreatedMoreThanOnce(0x3a2a6e13a488?)
	/Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__010/env/pkg/monitor/metrics_lifecycle_test.go:7 +0x24
testing.tRunner(0x3a2a6e13a488, 0x1013aee38)
	/opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2036 +0xc4
created by testing.(*T).Run in goroutine 1
	/opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2101 +0x3a8
FAIL	github.com/scheduler/pkg/monitor	0.646s
FAIL
```
