# Bug 复现说明

## Bug 是什么
`CronScheduler` 的任务表被多个 goroutine 并发访问：添加任务会写入任务表，移除任务会读取并删除任务表，但这些访问没有同步保护。因此 race detector 会报告数据竞争，运行时还可能因并发 map 操作终止。

## 如何触发
启动多个 goroutine，在统一起跑信号后反复交替执行任务添加和移除，并使用 Go race detector 运行测试。该测试包含 `WaitGroup` 和 channel 协调，能稳定触发共享任务表的并发访问。

## 运行指令
```bash
go test ./pkg/task -race -run '^TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree$' -count=1
```

## 错误信息
race detector 报告 `WARNING: DATA RACE`，访问路径涉及 `ScheduleTask` 与 `UnscheduleTask` 对同一任务表的读写。

## 错误堆栈
以下为上述命令在含 Bug 环境中的完整原始输出：
```text
==================
WARNING: DATA RACE
Write at 0x00c00011b500 by goroutine 11:
  runtime.mapaccess2_faststr()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/internal/runtime/maps/runtime_faststr.go:161 +0x2ac
  github.com/scheduler/pkg/task.(*CronScheduler).ScheduleTask()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron.go:37 +0xd0
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree.func1()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:31 +0x1d8

Previous write at 0x00c00011b500 by goroutine 14:
  runtime.mapassign_faststr()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/internal/runtime/maps/runtime_faststr.go:263 +0x4cc
  github.com/scheduler/pkg/task.(*CronScheduler).UnscheduleTask()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron.go:44 +0x2dc
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree.func1()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:34 +0x23c

Goroutine 11 (running) created at:
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:23 +0xc8
  testing.tRunner()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2036 +0x164
  testing.(*T).Run.gowrap1()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2101 +0x34

Goroutine 14 (running) created at:
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:23 +0xc8
  testing.tRunner()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2036 +0x164
  testing.(*T).Run.gowrap1()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2101 +0x34
==================
==================
WARNING: DATA RACE
Write at 0x00c000134c48 by goroutine 11:
  github.com/scheduler/pkg/task.(*CronScheduler).ScheduleTask()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron.go:37 +0xdc
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree.func1()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:31 +0x1d8

Previous write at 0x00c000134c48 by goroutine 14:
  github.com/scheduler/pkg/task.(*CronScheduler).ScheduleTask()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron.go:37 +0xdc
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree.func1()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:31 +0x1d8

Goroutine 11 (running) created at:
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:23 +0xc8
  testing.tRunner()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2036 +0x164
  testing.(*T).Run.gowrap1()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2101 +0x34

Goroutine 14 (running) created at:
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:23 +0xc8
  testing.tRunner()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2036 +0x164
  testing.(*T).Run.gowrap1()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2101 +0x34
==================
==================
WARNING: DATA RACE
Read at 0x00c00011b500 by goroutine 11:
  runtime.mapaccess1_faststr()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/internal/runtime/maps/runtime_faststr.go:101 +0x28c
  github.com/scheduler/pkg/task.(*CronScheduler).UnscheduleTask()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron.go:42 +0x27c
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree.func1()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:34 +0x23c

Previous write at 0x00c00011b500 by goroutine 12:
  runtime.mapassign_faststr()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/internal/runtime/maps/runtime_faststr.go:263 +0x4cc
  github.com/scheduler/pkg/task.(*CronScheduler).UnscheduleTask()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron.go:44 +0x2dc
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree.func1()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:34 +0x23c

Goroutine 11 (running) created at:
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:23 +0xc8
  testing.tRunner()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2036 +0x164
  testing.(*T).Run.gowrap1()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2101 +0x34

Goroutine 12 (running) created at:
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:23 +0xc8
  testing.tRunner()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2036 +0x164
  testing.(*T).Run.gowrap1()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2101 +0x34
==================
==================
WARNING: DATA RACE
Write at 0x00c00011b500 by goroutine 15:
  runtime.mapaccess2_faststr()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/internal/runtime/maps/runtime_faststr.go:161 +0x2ac
  github.com/scheduler/pkg/task.(*CronScheduler).ScheduleTask()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron.go:37 +0xd0
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree.func1()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:31 +0x1d8

Previous write at 0x00c00011b500 by goroutine 13:
  runtime.mapaccess2_faststr()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/internal/runtime/maps/runtime_faststr.go:161 +0x2ac
  github.com/scheduler/pkg/task.(*CronScheduler).ScheduleTask()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron.go:37 +0xd0
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree.func1()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:31 +0x1d8

Goroutine 15 (running) created at:
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:23 +0xc8
  testing.tRunner()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2036 +0x164
  testing.(*T).Run.gowrap1()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2101 +0x34

Goroutine 13 (running) created at:
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:23 +0xc8
  testing.tRunner()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2036 +0x164
  testing.(*T).Run.gowrap1()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2101 +0x34
==================
==================
WARNING: DATA RACE
Read at 0x00c00011b500 by goroutine 15:
  runtime.mapaccess1_faststr()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/internal/runtime/maps/runtime_faststr.go:101 +0x28c
  github.com/scheduler/pkg/task.(*CronScheduler).UnscheduleTask()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron.go:42 +0x27c
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree.func1()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:34 +0x23c

Previous write at 0x00c00011b500 by goroutine 9:
  runtime.mapaccess2_faststr()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/internal/runtime/maps/runtime_faststr.go:161 +0x2ac
  github.com/scheduler/pkg/task.(*CronScheduler).ScheduleTask()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron.go:37 +0xd0
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree.func1()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:31 +0x1d8

Goroutine 15 (running) created at:
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:23 +0xc8
  testing.tRunner()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2036 +0x164
  testing.(*T).Run.gowrap1()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2101 +0x34

Goroutine 9 (running) created at:
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:23 +0xc8
  testing.tRunner()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2036 +0x164
  testing.(*T).Run.gowrap1()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2101 +0x34
==================
==================
WARNING: DATA RACE
Write at 0x00c00011b500 by goroutine 8:
  runtime.mapassign_faststr()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/internal/runtime/maps/runtime_faststr.go:263 +0x4cc
  github.com/scheduler/pkg/task.(*CronScheduler).UnscheduleTask()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron.go:44 +0x2dc
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree.func1()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:34 +0x23c

Previous write at 0x00c00011b500 by goroutine 14:
  runtime.mapassign_faststr()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/internal/runtime/maps/runtime_faststr.go:263 +0x4cc
  github.com/scheduler/pkg/task.(*CronScheduler).UnscheduleTask()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron.go:44 +0x2dc
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree.func1()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:34 +0x23c

Goroutine 8 (running) created at:
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:23 +0xc8
  testing.tRunner()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2036 +0x164
  testing.(*T).Run.gowrap1()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2101 +0x34

Goroutine 14 (running) created at:
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:23 +0xc8
  testing.tRunner()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2036 +0x164
  testing.(*T).Run.gowrap1()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2101 +0x34
==================
==================
WARNING: DATA RACE
Write at 0x00c000134c90 by goroutine 14:
  github.com/scheduler/pkg/task.(*CronScheduler).ScheduleTask()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron.go:37 +0xdc
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree.func1()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:31 +0x1d8

Previous read at 0x00c000134c90 by goroutine 15:
  github.com/scheduler/pkg/task.(*CronScheduler).UnscheduleTask()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron.go:42 +0x288
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree.func1()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:34 +0x23c

Goroutine 14 (running) created at:
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:23 +0xc8
  testing.tRunner()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2036 +0x164
  testing.(*T).Run.gowrap1()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2101 +0x34

Goroutine 15 (running) created at:
  github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree()
      /Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:23 +0xc8
  testing.tRunner()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2036 +0x164
  testing.(*T).Run.gowrap1()
      /opt/homebrew/Cellar/go/1.26.5/libexec/src/testing/testing.go:2101 +0x34
==================
fatal error: concurrent map read and map write

goroutine 25 [running]:
internal/runtime/maps.fatal({0x100490413?, 0x100472abc?})
	/opt/homebrew/Cellar/go/1.26.5/libexec/src/runtime/panic.go:1181 +0x20
github.com/scheduler/pkg/task.(*CronScheduler).UnscheduleTask(...)
	/Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron.go:42
github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree.func1()
	/Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:34 +0x280
created by github.com/scheduler/pkg/task.TestCronSchedulerConcurrentScheduleAndUnscheduleIsRaceFree in goroutine 19
	/Users/tog/Desktop/code/go标注/我的go/2026-08-19/cdn-product__009/env/pkg/task/cron_concurrency_test.go:23 +0xcc
FAIL	github.com/scheduler/pkg/task	0.559s
FAIL

```
