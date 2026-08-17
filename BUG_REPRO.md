# 缺陷复现说明

## 缺陷现象
通过 API 创建不带 cron 的 cache-warm 任务后，接口返回 201、数据库也能查到，但 scheduler 永远不执行它。手动任务也应该进入就绪队列，请修复这个问题。

## 触发方式
# 1. 经 API 注册一个没有 cron 表达式的任务，检查调度 DAG 是否可返回该就绪任务。
go test -v ./pkg/scheduler -run '^TestRegisterTaskKeepsManualTaskReadyForScheduler$' -count=20
# 2. 回归调度器及其依赖包的全部测试。
go test ./...

## 触发后的实际错误输出

```text
=== RUN   TestRegisterTaskKeepsManualTaskReadyForScheduler
    task_registration_test.go:27: manual API task was not available to the scheduler: []*task.Task(nil)
--- FAIL: TestRegisterTaskKeepsManualTaskReadyForScheduler (0.00s)
=== RUN   TestRegisterTaskKeepsManualTaskReadyForScheduler
    task_registration_test.go:27: manual API task was not available to the scheduler: []*task.Task(nil)
--- FAIL: TestRegisterTaskKeepsManualTaskReadyForScheduler (0.00s)
=== RUN   TestRegisterTaskKeepsManualTaskReadyForScheduler
    task_registration_test.go:27: manual API task was not available to the scheduler: []*task.Task(nil)
--- FAIL: TestRegisterTaskKeepsManualTaskReadyForScheduler (0.00s)
=== RUN   TestRegisterTaskKeepsManualTaskReadyForScheduler
    task_registration_test.go:27: manual API task was not available to the scheduler: []*task.Task(nil)
--- FAIL: TestRegisterTaskKeepsManualTaskReadyForScheduler (0.00s)
=== RUN   TestRegisterTaskKeepsManualTaskReadyForScheduler
    task_registration_test.go:27: manual API task was not available to the scheduler: []*task.Task(nil)
--- FAIL: TestRegisterTaskKeepsManualTaskReadyForScheduler (0.00s)
=== RUN   TestRegisterTaskKeepsManualTaskReadyForScheduler
    task_registration_test.go:27: manual API task was not available to the scheduler: []*task.Task(nil)
--- FAIL: TestRegisterTaskKeepsManualTaskReadyForScheduler (0.00s)
=== RUN   TestRegisterTaskKeepsManualTaskReadyForScheduler
    task_registration_test.go:27: manual API task was not available to the scheduler: []*task.Task(nil)
--- FAIL: TestRegisterTaskKeepsManualTaskReadyForScheduler (0.00s)
=== RUN   TestRegisterTaskKeepsManualTaskReadyForScheduler
    task_registration_test.go:27: manual API task was not available to the scheduler: []*task.Task(nil)
--- FAIL: TestRegisterTaskKeepsManualTaskReadyForScheduler (0.00s)
=== RUN   TestRegisterTaskKeepsManualTaskReadyForScheduler
    task_registration_test.go:27: manual API task was not available to the scheduler: []*task.Task(nil)
--- FAIL: TestRegisterTaskKeepsManualTaskReadyForScheduler (0.00s)
=== RUN   TestRegisterTaskKeepsManualTaskReadyForScheduler
    task_registration_test.go:27: manual API task was not available to the scheduler: []*task.Task(nil)
--- FAIL: TestRegisterTaskKeepsManualTaskReadyForScheduler (0.00s)
=== RUN   TestRegisterTaskKeepsManualTaskReadyForScheduler
    task_registration_test.go:27: manual API task was not available to the scheduler: []*task.Task(nil)
--- FAIL: TestRegisterTaskKeepsManualTaskReadyForScheduler (0.00s)
=== RUN   TestRegisterTaskKeepsManualTaskReadyForScheduler
    task_registration_test.go:27: manual API task was not available to the scheduler: []*task.Task(nil)
--- FAIL: TestRegisterTaskKeepsManualTaskReadyForScheduler (0.00s)
=== RUN   TestRegisterTaskKeepsManualTaskReadyForScheduler
    task_registration_test.go:27: manual API task was not available to the scheduler: []*task.Task(nil)
--- FAIL: TestRegisterTaskKeepsManualTaskReadyForScheduler (0.00s)
=== RUN   TestRegisterTaskKeepsManualTaskReadyForScheduler
    task_registration_test.go:27: manual API task was not available to the scheduler: []*task.Task(nil)
--- FAIL: TestRegisterTaskKeepsManualTaskReadyForScheduler (0.00s)
=== RUN   TestRegisterTaskKeepsManualTaskReadyForScheduler
    task_registration_test.go:27: manual API task was not available to the scheduler: []*task.Task(nil)
--- FAIL: TestRegisterTaskKeepsManualTaskReadyForScheduler (0.00s)
=== RUN   TestRegisterTaskKeepsManualTaskReadyForScheduler
    task_registration_test.go:27: manual API task was not available to the scheduler: []*task.Task(nil)
--- FAIL: TestRegisterTaskKeepsManualTaskReadyForScheduler (0.00s)
=== RUN   TestRegisterTaskKeepsManualTaskReadyForScheduler
    task_registration_test.go:27: manual API task was not available to the scheduler: []*task.Task(nil)
--- FAIL: TestRegisterTaskKeepsManualTaskReadyForScheduler (0.00s)
=== RUN   TestRegisterTaskKeepsManualTaskReadyForScheduler
    task_registration_test.go:27: manual API task was not available to the scheduler: []*task.Task(nil)
--- FAIL: TestRegisterTaskKeepsManualTaskReadyForScheduler (0.00s)
=== RUN   TestRegisterTaskKeepsManualTaskReadyForScheduler
    task_registration_test.go:27: manual API task was not available to the scheduler: []*task.Task(nil)
--- FAIL: TestRegisterTaskKeepsManualTaskReadyForScheduler (0.00s)
=== RUN   TestRegisterTaskKeepsManualTaskReadyForScheduler
    task_registration_test.go:27: manual API task was not available to the scheduler: []*task.Task(nil)
--- FAIL: TestRegisterTaskKeepsManualTaskReadyForScheduler (0.00s)
FAIL
FAIL	github.com/scheduler/pkg/scheduler	0.006s
FAIL
```
