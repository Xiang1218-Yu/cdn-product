# 缺陷复现说明

## 缺陷现象
我们的 scheduler 给 Capacity=1 的 executor 派了一个任务后，同一轮又把另一个任务派过去，下游因此出现并发写入冲突。容量应该在任务结束前一直被占用，请修复这个问题。

## 触发方式
# 1. 容量为 1 时，第二个不同任务必须在第一个任务释放前被拒绝。
go test -race -v ./pkg/scheduler -run '^TestAcquireExecutorReservesCapacityUntilTaskCompletes$' -count=20

## 触发后的实际错误输出

```text
=== RUN   TestAcquireExecutorReservesCapacityUntilTaskCompletes
    executor_capacity_test.go:14: capacity=1 admitted a second task before the first completed: &scheduler.Executor{ID:"one-slot", Address:"", Capacity:1, Load:0, LastSeen:time.Date(2026, time.August, 17, 1, 15, 20, 112980964, time.Local)}
--- FAIL: TestAcquireExecutorReservesCapacityUntilTaskCompletes (0.00s)
=== RUN   TestAcquireExecutorReservesCapacityUntilTaskCompletes
    executor_capacity_test.go:14: capacity=1 admitted a second task before the first completed: &scheduler.Executor{ID:"one-slot", Address:"", Capacity:1, Load:0, LastSeen:time.Date(2026, time.August, 17, 1, 15, 20, 113168589, time.Local)}
--- FAIL: TestAcquireExecutorReservesCapacityUntilTaskCompletes (0.00s)
=== RUN   TestAcquireExecutorReservesCapacityUntilTaskCompletes
    executor_capacity_test.go:14: capacity=1 admitted a second task before the first completed: &scheduler.Executor{ID:"one-slot", Address:"", Capacity:1, Load:0, LastSeen:time.Date(2026, time.August, 17, 1, 15, 20, 113320380, time.Local)}
--- FAIL: TestAcquireExecutorReservesCapacityUntilTaskCompletes (0.00s)
=== RUN   TestAcquireExecutorReservesCapacityUntilTaskCompletes
    executor_capacity_test.go:14: capacity=1 admitted a second task before the first completed: &scheduler.Executor{ID:"one-slot", Address:"", Capacity:1, Load:0, LastSeen:time.Date(2026, time.August, 17, 1, 15, 20, 113457797, time.Local)}
--- FAIL: TestAcquireExecutorReservesCapacityUntilTaskCompletes (0.00s)
=== RUN   TestAcquireExecutorReservesCapacityUntilTaskCompletes
    executor_capacity_test.go:14: capacity=1 admitted a second task before the first completed: &scheduler.Executor{ID:"one-slot", Address:"", Capacity:1, Load:0, LastSeen:time.Date(2026, time.August, 17, 1, 15, 20, 113734630, time.Local)}
--- FAIL: TestAcquireExecutorReservesCapacityUntilTaskCompletes (0.00s)
=== RUN   TestAcquireExecutorReservesCapacityUntilTaskCompletes
    executor_capacity_test.go:14: capacity=1 admitted a second task before the first completed: &scheduler.Executor{ID:"one-slot", Address:"", Capacity:1, Load:0, LastSeen:time.Date(2026, time.August, 17, 1, 15, 20, 113844547, time.Local)}
--- FAIL: TestAcquireExecutorReservesCapacityUntilTaskCompletes (0.00s)
=== RUN   TestAcquireExecutorReservesCapacityUntilTaskCompletes
    executor_capacity_test.go:14: capacity=1 admitted a second task before the first completed: &scheduler.Executor{ID:"one-slot", Address:"", Capacity:1, Load:0, LastSeen:time.Date(2026, time.August, 17, 1, 15, 20, 113949589, time.Local)}
--- FAIL: TestAcquireExecutorReservesCapacityUntilTaskCompletes (0.00s)
=== RUN   TestAcquireExecutorReservesCapacityUntilTaskCompletes
    executor_capacity_test.go:14: capacity=1 admitted a second task before the first completed: &scheduler.Executor{ID:"one-slot", Address:"", Capacity:1, Load:0, LastSeen:time.Date(2026, time.August, 17, 1, 15, 20, 114053547, time.Local)}
--- FAIL: TestAcquireExecutorReservesCapacityUntilTaskCompletes (0.00s)
=== RUN   TestAcquireExecutorReservesCapacityUntilTaskCompletes
    executor_capacity_test.go:14: capacity=1 admitted a second task before the first completed: &scheduler.Executor{ID:"one-slot", Address:"", Capacity:1, Load:0, LastSeen:time.Date(2026, time.August, 17, 1, 15, 20, 114177172, time.Local)}
--- FAIL: TestAcquireExecutorReservesCapacityUntilTaskCompletes (0.00s)
=== RUN   TestAcquireExecutorReservesCapacityUntilTaskCompletes
    executor_capacity_test.go:14: capacity=1 admitted a second task before the first completed: &scheduler.Executor{ID:"one-slot", Address:"", Capacity:1, Load:0, LastSeen:time.Date(2026, time.August, 17, 1, 15, 20, 114293047, time.Local)}
--- FAIL: TestAcquireExecutorReservesCapacityUntilTaskCompletes (0.00s)
=== RUN   TestAcquireExecutorReservesCapacityUntilTaskCompletes
    executor_capacity_test.go:14: capacity=1 admitted a second task before the first completed: &scheduler.Executor{ID:"one-slot", Address:"", Capacity:1, Load:0, LastSeen:time.Date(2026, time.August, 17, 1, 15, 20, 114414714, time.Local)}
--- FAIL: TestAcquireExecutorReservesCapacityUntilTaskCompletes (0.00s)
=== RUN   TestAcquireExecutorReservesCapacityUntilTaskCompletes
    executor_capacity_test.go:14: capacity=1 admitted a second task before the first completed: &scheduler.Executor{ID:"one-slot", Address:"", Capacity:1, Load:0, LastSeen:time.Date(2026, time.August, 17, 1, 15, 20, 114623547, time.Local)}
--- FAIL: TestAcquireExecutorReservesCapacityUntilTaskCompletes (0.00s)
=== RUN   TestAcquireExecutorReservesCapacityUntilTaskCompletes
    executor_capacity_test.go:14: capacity=1 admitted a second task before the first completed: &scheduler.Executor{ID:"one-slot", Address:"", Capacity:1, Load:0, LastSeen:time.Date(2026, time.August, 17, 1, 15, 20, 114756922, time.Local)}
--- FAIL: TestAcquireExecutorReservesCapacityUntilTaskCompletes (0.00s)
=== RUN   TestAcquireExecutorReservesCapacityUntilTaskCompletes
    executor_capacity_test.go:14: capacity=1 admitted a second task before the first completed: &scheduler.Executor{ID:"one-slot", Address:"", Capacity:1, Load:0, LastSeen:time.Date(2026, time.August, 17, 1, 15, 20, 114850130, time.Local)}
--- FAIL: TestAcquireExecutorReservesCapacityUntilTaskCompletes (0.00s)
=== RUN   TestAcquireExecutorReservesCapacityUntilTaskCompletes
    executor_capacity_test.go:14: capacity=1 admitted a second task before the first completed: &scheduler.Executor{ID:"one-slot", Address:"", Capacity:1, Load:0, LastSeen:time.Date(2026, time.August, 17, 1, 15, 20, 114942339, time.Local)}
--- FAIL: TestAcquireExecutorReservesCapacityUntilTaskCompletes (0.00s)
=== RUN   TestAcquireExecutorReservesCapacityUntilTaskCompletes
    executor_capacity_test.go:14: capacity=1 admitted a second task before the first completed: &scheduler.Executor{ID:"one-slot", Address:"", Capacity:1, Load:0, LastSeen:time.Date(2026, time.August, 17, 1, 15, 20, 115145047, time.Local)}
--- FAIL: TestAcquireExecutorReservesCapacityUntilTaskCompletes (0.00s)
=== RUN   TestAcquireExecutorReservesCapacityUntilTaskCompletes
    executor_capacity_test.go:14: capacity=1 admitted a second task before the first completed: &scheduler.Executor{ID:"one-slot", Address:"", Capacity:1, Load:0, LastSeen:time.Date(2026, time.August, 17, 1, 15, 20, 115283089, time.Local)}
--- FAIL: TestAcquireExecutorReservesCapacityUntilTaskCompletes (0.00s)
=== RUN   TestAcquireExecutorReservesCapacityUntilTaskCompletes
    executor_capacity_test.go:14: capacity=1 admitted a second task before the first completed: &scheduler.Executor{ID:"one-slot", Address:"", Capacity:1, Load:0, LastSeen:time.Date(2026, time.August, 17, 1, 15, 20, 115378547, time.Local)}
--- FAIL: TestAcquireExecutorReservesCapacityUntilTaskCompletes (0.00s)
=== RUN   TestAcquireExecutorReservesCapacityUntilTaskCompletes
    executor_capacity_test.go:14: capacity=1 admitted a second task before the first completed: &scheduler.Executor{ID:"one-slot", Address:"", Capacity:1, Load:0, LastSeen:time.Date(2026, time.August, 17, 1, 15, 20, 115485130, time.Local)}
--- FAIL: TestAcquireExecutorReservesCapacityUntilTaskCompletes (0.00s)
=== RUN   TestAcquireExecutorReservesCapacityUntilTaskCompletes
    executor_capacity_test.go:14: capacity=1 admitted a second task before the first completed: &scheduler.Executor{ID:"one-slot", Address:"", Capacity:1, Load:0, LastSeen:time.Date(2026, time.August, 17, 1, 15, 20, 115575839, time.Local)}
--- FAIL: TestAcquireExecutorReservesCapacityUntilTaskCompletes (0.00s)
FAIL
FAIL	github.com/scheduler/pkg/scheduler	0.016s
FAIL
```
