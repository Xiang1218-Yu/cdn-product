# 缺陷复现说明

## 缺陷现象
executor dead 后日志已经选到了 survivor，但 report 的故障记录还是指向 dead；下一次巡检又重复处理它，任务一直卡着。先不要改文件，帮我查清这个 failover 为什么没有真正接管执行记录。

## 触发方式
# 1. 记录由 dead 执行的 report，提供 survivor 后触发 failover 并断言执行记录已切换。
go test -v ./pkg/scheduler -run '^TestFailoverRecordsReplacementExecutor$' -count=20

## 触发后的实际错误输出

```text
=== RUN   TestFailoverRecordsReplacementExecutor
    failover_reassignment_test.go:32: failover selected survivor but kept stale execution record: &scheduler.TaskExecution{TaskID:"report", ExecutorID:"dead", StartTime:time.Date(2026, time.August, 17, 10, 30, 39, 565964000, time.Local), LastUpdate:time.Date(2026, time.August, 17, 10, 30, 39, 565964000, time.Local)}
--- FAIL: TestFailoverRecordsReplacementExecutor (0.00s)
=== RUN   TestFailoverRecordsReplacementExecutor
    failover_reassignment_test.go:32: failover selected survivor but kept stale execution record: &scheduler.TaskExecution{TaskID:"report", ExecutorID:"dead", StartTime:time.Date(2026, time.August, 17, 10, 30, 39, 566406000, time.Local), LastUpdate:time.Date(2026, time.August, 17, 10, 30, 39, 566406000, time.Local)}
--- FAIL: TestFailoverRecordsReplacementExecutor (0.00s)
=== RUN   TestFailoverRecordsReplacementExecutor
    failover_reassignment_test.go:32: failover selected survivor but kept stale execution record: &scheduler.TaskExecution{TaskID:"report", ExecutorID:"dead", StartTime:time.Date(2026, time.August, 17, 10, 30, 39, 566454000, time.Local), LastUpdate:time.Date(2026, time.August, 17, 10, 30, 39, 566454000, time.Local)}
--- FAIL: TestFailoverRecordsReplacementExecutor (0.00s)
=== RUN   TestFailoverRecordsReplacementExecutor
    failover_reassignment_test.go:32: failover selected survivor but kept stale execution record: &scheduler.TaskExecution{TaskID:"report", ExecutorID:"dead", StartTime:time.Date(2026, time.August, 17, 10, 30, 39, 566491000, time.Local), LastUpdate:time.Date(2026, time.August, 17, 10, 30, 39, 566491000, time.Local)}
--- FAIL: TestFailoverRecordsReplacementExecutor (0.00s)
=== RUN   TestFailoverRecordsReplacementExecutor
    failover_reassignment_test.go:32: failover selected survivor but kept stale execution record: &scheduler.TaskExecution{TaskID:"report", ExecutorID:"dead", StartTime:time.Date(2026, time.August, 17, 10, 30, 39, 566520000, time.Local), LastUpdate:time.Date(2026, time.August, 17, 10, 30, 39, 566520000, time.Local)}
--- FAIL: TestFailoverRecordsReplacementExecutor (0.00s)
=== RUN   TestFailoverRecordsReplacementExecutor
    failover_reassignment_test.go:32: failover selected survivor but kept stale execution record: &scheduler.TaskExecution{TaskID:"report", ExecutorID:"dead", StartTime:time.Date(2026, time.August, 17, 10, 30, 39, 566715000, time.Local), LastUpdate:time.Date(2026, time.August, 17, 10, 30, 39, 566715000, time.Local)}
--- FAIL: TestFailoverRecordsReplacementExecutor (0.00s)
=== RUN   TestFailoverRecordsReplacementExecutor
    failover_reassignment_test.go:32: failover selected survivor but kept stale execution record: &scheduler.TaskExecution{TaskID:"report", ExecutorID:"dead", StartTime:time.Date(2026, time.August, 17, 10, 30, 39, 566746000, time.Local), LastUpdate:time.Date(2026, time.August, 17, 10, 30, 39, 566746000, time.Local)}
--- FAIL: TestFailoverRecordsReplacementExecutor (0.00s)
=== RUN   TestFailoverRecordsReplacementExecutor
    failover_reassignment_test.go:32: failover selected survivor but kept stale execution record: &scheduler.TaskExecution{TaskID:"report", ExecutorID:"dead", StartTime:time.Date(2026, time.August, 17, 10, 30, 39, 566776000, time.Local), LastUpdate:time.Date(2026, time.August, 17, 10, 30, 39, 566776000, time.Local)}
--- FAIL: TestFailoverRecordsReplacementExecutor (0.00s)
=== RUN   TestFailoverRecordsReplacementExecutor
    failover_reassignment_test.go:32: failover selected survivor but kept stale execution record: &scheduler.TaskExecution{TaskID:"report", ExecutorID:"dead", StartTime:time.Date(2026, time.August, 17, 10, 30, 39, 566807000, time.Local), LastUpdate:time.Date(2026, time.August, 17, 10, 30, 39, 566807000, time.Local)}
--- FAIL: TestFailoverRecordsReplacementExecutor (0.00s)
=== RUN   TestFailoverRecordsReplacementExecutor
    failover_reassignment_test.go:32: failover selected survivor but kept stale execution record: &scheduler.TaskExecution{TaskID:"report", ExecutorID:"dead", StartTime:time.Date(2026, time.August, 17, 10, 30, 39, 566836000, time.Local), LastUpdate:time.Date(2026, time.August, 17, 10, 30, 39, 566836000, time.Local)}
--- FAIL: TestFailoverRecordsReplacementExecutor (0.00s)
=== RUN   TestFailoverRecordsReplacementExecutor
    failover_reassignment_test.go:32: failover selected survivor but kept stale execution record: &scheduler.TaskExecution{TaskID:"report", ExecutorID:"dead", StartTime:time.Date(2026, time.August, 17, 10, 30, 39, 566874000, time.Local), LastUpdate:time.Date(2026, time.August, 17, 10, 30, 39, 566874000, time.Local)}
--- FAIL: TestFailoverRecordsReplacementExecutor (0.00s)
=== RUN   TestFailoverRecordsReplacementExecutor
    failover_reassignment_test.go:32: failover selected survivor but kept stale execution record: &scheduler.TaskExecution{TaskID:"report", ExecutorID:"dead", StartTime:time.Date(2026, time.August, 17, 10, 30, 39, 566907000, time.Local), LastUpdate:time.Date(2026, time.August, 17, 10, 30, 39, 566907000, time.Local)}
--- FAIL: TestFailoverRecordsReplacementExecutor (0.00s)
=== RUN   TestFailoverRecordsReplacementExecutor
    failover_reassignment_test.go:32: failover selected survivor but kept stale execution record: &scheduler.TaskExecution{TaskID:"report", ExecutorID:"dead", StartTime:time.Date(2026, time.August, 17, 10, 30, 39, 566940000, time.Local), LastUpdate:time.Date(2026, time.August, 17, 10, 30, 39, 566940000, time.Local)}
--- FAIL: TestFailoverRecordsReplacementExecutor (0.00s)
=== RUN   TestFailoverRecordsReplacementExecutor
    failover_reassignment_test.go:32: failover selected survivor but kept stale execution record: &scheduler.TaskExecution{TaskID:"report", ExecutorID:"dead", StartTime:time.Date(2026, time.August, 17, 10, 30, 39, 566978000, time.Local), LastUpdate:time.Date(2026, time.August, 17, 10, 30, 39, 566978000, time.Local)}
--- FAIL: TestFailoverRecordsReplacementExecutor (0.00s)
=== RUN   TestFailoverRecordsReplacementExecutor
    failover_reassignment_test.go:32: failover selected survivor but kept stale execution record: &scheduler.TaskExecution{TaskID:"report", ExecutorID:"dead", StartTime:time.Date(2026, time.August, 17, 10, 30, 39, 567011000, time.Local), LastUpdate:time.Date(2026, time.August, 17, 10, 30, 39, 567011000, time.Local)}
--- FAIL: TestFailoverRecordsReplacementExecutor (0.00s)
=== RUN   TestFailoverRecordsReplacementExecutor
    failover_reassignment_test.go:32: failover selected survivor but kept stale execution record: &scheduler.TaskExecution{TaskID:"report", ExecutorID:"dead", StartTime:time.Date(2026, time.August, 17, 10, 30, 39, 567040000, time.Local), LastUpdate:time.Date(2026, time.August, 17, 10, 30, 39, 567041000, time.Local)}
--- FAIL: TestFailoverRecordsReplacementExecutor (0.00s)
=== RUN   TestFailoverRecordsReplacementExecutor
    failover_reassignment_test.go:32: failover selected survivor but kept stale execution record: &scheduler.TaskExecution{TaskID:"report", ExecutorID:"dead", StartTime:time.Date(2026, time.August, 17, 10, 30, 39, 567074000, time.Local), LastUpdate:time.Date(2026, time.August, 17, 10, 30, 39, 567074000, time.Local)}
--- FAIL: TestFailoverRecordsReplacementExecutor (0.00s)
=== RUN   TestFailoverRecordsReplacementExecutor
    failover_reassignment_test.go:32: failover selected survivor but kept stale execution record: &scheduler.TaskExecution{TaskID:"report", ExecutorID:"dead", StartTime:time.Date(2026, time.August, 17, 10, 30, 39, 567103000, time.Local), LastUpdate:time.Date(2026, time.August, 17, 10, 30, 39, 567103000, time.Local)}
--- FAIL: TestFailoverRecordsReplacementExecutor (0.00s)
=== RUN   TestFailoverRecordsReplacementExecutor
    failover_reassignment_test.go:32: failover selected survivor but kept stale execution record: &scheduler.TaskExecution{TaskID:"report", ExecutorID:"dead", StartTime:time.Date(2026, time.August, 17, 10, 30, 39, 567133000, time.Local), LastUpdate:time.Date(2026, time.August, 17, 10, 30, 39, 567133000, time.Local)}
--- FAIL: TestFailoverRecordsReplacementExecutor (0.00s)
=== RUN   TestFailoverRecordsReplacementExecutor
    failover_reassignment_test.go:32: failover selected survivor but kept stale execution record: &scheduler.TaskExecution{TaskID:"report", ExecutorID:"dead", StartTime:time.Date(2026, time.August, 17, 10, 30, 39, 567220000, time.Local), LastUpdate:time.Date(2026, time.August, 17, 10, 30, 39, 567220000, time.Local)}
--- FAIL: TestFailoverRecordsReplacementExecutor (0.00s)
FAIL
FAIL	github.com/scheduler/pkg/scheduler	1.114s
FAIL
```
