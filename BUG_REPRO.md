# 缺陷复现说明

## 缺陷现象
编辑已经存在的定时任务后，daily-report 会按旧表达式和新表达式各跑一次，重复发送日报。先别改代码，帮我定位为什么同一个任务 ID 会留下两条调度记录。

## 触发方式
# 1. 以同一任务 ID 连续注册两种 cron 表达式，并检查底层活动 entry 的数量。
go test -v ./pkg/task -run '^TestReschedulingSameTaskReplacesPreviousCronEntry$' -count=20

## 触发后的实际错误输出

```text
=== RUN   TestReschedulingSameTaskReplacesPreviousCronEntry
    cron_replacement_test.go:21: rescheduling one task left 2 active cron entries, want 1
--- FAIL: TestReschedulingSameTaskReplacesPreviousCronEntry (0.00s)
=== RUN   TestReschedulingSameTaskReplacesPreviousCronEntry
    cron_replacement_test.go:21: rescheduling one task left 2 active cron entries, want 1
--- FAIL: TestReschedulingSameTaskReplacesPreviousCronEntry (0.00s)
=== RUN   TestReschedulingSameTaskReplacesPreviousCronEntry
    cron_replacement_test.go:21: rescheduling one task left 2 active cron entries, want 1
--- FAIL: TestReschedulingSameTaskReplacesPreviousCronEntry (0.00s)
=== RUN   TestReschedulingSameTaskReplacesPreviousCronEntry
    cron_replacement_test.go:21: rescheduling one task left 2 active cron entries, want 1
--- FAIL: TestReschedulingSameTaskReplacesPreviousCronEntry (0.00s)
=== RUN   TestReschedulingSameTaskReplacesPreviousCronEntry
    cron_replacement_test.go:21: rescheduling one task left 2 active cron entries, want 1
--- FAIL: TestReschedulingSameTaskReplacesPreviousCronEntry (0.00s)
=== RUN   TestReschedulingSameTaskReplacesPreviousCronEntry
    cron_replacement_test.go:21: rescheduling one task left 2 active cron entries, want 1
--- FAIL: TestReschedulingSameTaskReplacesPreviousCronEntry (0.00s)
=== RUN   TestReschedulingSameTaskReplacesPreviousCronEntry
    cron_replacement_test.go:21: rescheduling one task left 2 active cron entries, want 1
--- FAIL: TestReschedulingSameTaskReplacesPreviousCronEntry (0.00s)
=== RUN   TestReschedulingSameTaskReplacesPreviousCronEntry
    cron_replacement_test.go:21: rescheduling one task left 2 active cron entries, want 1
--- FAIL: TestReschedulingSameTaskReplacesPreviousCronEntry (0.00s)
=== RUN   TestReschedulingSameTaskReplacesPreviousCronEntry
    cron_replacement_test.go:21: rescheduling one task left 2 active cron entries, want 1
--- FAIL: TestReschedulingSameTaskReplacesPreviousCronEntry (0.00s)
=== RUN   TestReschedulingSameTaskReplacesPreviousCronEntry
    cron_replacement_test.go:21: rescheduling one task left 2 active cron entries, want 1
--- FAIL: TestReschedulingSameTaskReplacesPreviousCronEntry (0.00s)
=== RUN   TestReschedulingSameTaskReplacesPreviousCronEntry
    cron_replacement_test.go:21: rescheduling one task left 2 active cron entries, want 1
--- FAIL: TestReschedulingSameTaskReplacesPreviousCronEntry (0.00s)
=== RUN   TestReschedulingSameTaskReplacesPreviousCronEntry
    cron_replacement_test.go:21: rescheduling one task left 2 active cron entries, want 1
--- FAIL: TestReschedulingSameTaskReplacesPreviousCronEntry (0.00s)
=== RUN   TestReschedulingSameTaskReplacesPreviousCronEntry
    cron_replacement_test.go:21: rescheduling one task left 2 active cron entries, want 1
--- FAIL: TestReschedulingSameTaskReplacesPreviousCronEntry (0.00s)
=== RUN   TestReschedulingSameTaskReplacesPreviousCronEntry
    cron_replacement_test.go:21: rescheduling one task left 2 active cron entries, want 1
--- FAIL: TestReschedulingSameTaskReplacesPreviousCronEntry (0.00s)
=== RUN   TestReschedulingSameTaskReplacesPreviousCronEntry
    cron_replacement_test.go:21: rescheduling one task left 2 active cron entries, want 1
--- FAIL: TestReschedulingSameTaskReplacesPreviousCronEntry (0.00s)
=== RUN   TestReschedulingSameTaskReplacesPreviousCronEntry
    cron_replacement_test.go:21: rescheduling one task left 2 active cron entries, want 1
--- FAIL: TestReschedulingSameTaskReplacesPreviousCronEntry (0.00s)
=== RUN   TestReschedulingSameTaskReplacesPreviousCronEntry
    cron_replacement_test.go:21: rescheduling one task left 2 active cron entries, want 1
--- FAIL: TestReschedulingSameTaskReplacesPreviousCronEntry (0.00s)
=== RUN   TestReschedulingSameTaskReplacesPreviousCronEntry
    cron_replacement_test.go:21: rescheduling one task left 2 active cron entries, want 1
--- FAIL: TestReschedulingSameTaskReplacesPreviousCronEntry (0.00s)
=== RUN   TestReschedulingSameTaskReplacesPreviousCronEntry
    cron_replacement_test.go:21: rescheduling one task left 2 active cron entries, want 1
--- FAIL: TestReschedulingSameTaskReplacesPreviousCronEntry (0.00s)
=== RUN   TestReschedulingSameTaskReplacesPreviousCronEntry
    cron_replacement_test.go:21: rescheduling one task left 2 active cron entries, want 1
--- FAIL: TestReschedulingSameTaskReplacesPreviousCronEntry (0.00s)
FAIL
FAIL	github.com/scheduler/pkg/task	0.478s
FAIL
```
