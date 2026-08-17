# 缺陷复现说明

## 缺陷现象
我们给已经存在的 worker 更新依赖，更新因为环依赖被拒绝后，原来的 worker 反而从调度图里消失，后续 report 也一直不跑。文件先不要改，帮我定位这个问题的根因。

## 触发方式
# 1. 构造被拒绝的环依赖替换，并检查原 worker 是否还可读取和调度。
go test -v ./pkg/task -run '^TestRejectedCyclicReplacementKeepsExistingTask$' -count=20

## 触发后的实际错误输出

```text
=== RUN   TestRejectedCyclicReplacementKeepsExistingTask
    dag_replacement_test.go:32: rejected update removed the existing worker: task not found
--- FAIL: TestRejectedCyclicReplacementKeepsExistingTask (0.00s)
=== RUN   TestRejectedCyclicReplacementKeepsExistingTask
    dag_replacement_test.go:32: rejected update removed the existing worker: task not found
--- FAIL: TestRejectedCyclicReplacementKeepsExistingTask (0.00s)
=== RUN   TestRejectedCyclicReplacementKeepsExistingTask
    dag_replacement_test.go:32: rejected update removed the existing worker: task not found
--- FAIL: TestRejectedCyclicReplacementKeepsExistingTask (0.00s)
=== RUN   TestRejectedCyclicReplacementKeepsExistingTask
    dag_replacement_test.go:32: rejected update removed the existing worker: task not found
--- FAIL: TestRejectedCyclicReplacementKeepsExistingTask (0.00s)
=== RUN   TestRejectedCyclicReplacementKeepsExistingTask
    dag_replacement_test.go:32: rejected update removed the existing worker: task not found
--- FAIL: TestRejectedCyclicReplacementKeepsExistingTask (0.00s)
=== RUN   TestRejectedCyclicReplacementKeepsExistingTask
    dag_replacement_test.go:32: rejected update removed the existing worker: task not found
--- FAIL: TestRejectedCyclicReplacementKeepsExistingTask (0.00s)
=== RUN   TestRejectedCyclicReplacementKeepsExistingTask
    dag_replacement_test.go:32: rejected update removed the existing worker: task not found
--- FAIL: TestRejectedCyclicReplacementKeepsExistingTask (0.00s)
=== RUN   TestRejectedCyclicReplacementKeepsExistingTask
    dag_replacement_test.go:32: rejected update removed the existing worker: task not found
--- FAIL: TestRejectedCyclicReplacementKeepsExistingTask (0.00s)
=== RUN   TestRejectedCyclicReplacementKeepsExistingTask
    dag_replacement_test.go:32: rejected update removed the existing worker: task not found
--- FAIL: TestRejectedCyclicReplacementKeepsExistingTask (0.00s)
=== RUN   TestRejectedCyclicReplacementKeepsExistingTask
    dag_replacement_test.go:32: rejected update removed the existing worker: task not found
--- FAIL: TestRejectedCyclicReplacementKeepsExistingTask (0.00s)
=== RUN   TestRejectedCyclicReplacementKeepsExistingTask
    dag_replacement_test.go:32: rejected update removed the existing worker: task not found
--- FAIL: TestRejectedCyclicReplacementKeepsExistingTask (0.00s)
=== RUN   TestRejectedCyclicReplacementKeepsExistingTask
    dag_replacement_test.go:32: rejected update removed the existing worker: task not found
--- FAIL: TestRejectedCyclicReplacementKeepsExistingTask (0.00s)
=== RUN   TestRejectedCyclicReplacementKeepsExistingTask
    dag_replacement_test.go:32: rejected update removed the existing worker: task not found
--- FAIL: TestRejectedCyclicReplacementKeepsExistingTask (0.00s)
=== RUN   TestRejectedCyclicReplacementKeepsExistingTask
    dag_replacement_test.go:32: rejected update removed the existing worker: task not found
--- FAIL: TestRejectedCyclicReplacementKeepsExistingTask (0.00s)
=== RUN   TestRejectedCyclicReplacementKeepsExistingTask
    dag_replacement_test.go:32: rejected update removed the existing worker: task not found
--- FAIL: TestRejectedCyclicReplacementKeepsExistingTask (0.00s)
=== RUN   TestRejectedCyclicReplacementKeepsExistingTask
    dag_replacement_test.go:32: rejected update removed the existing worker: task not found
--- FAIL: TestRejectedCyclicReplacementKeepsExistingTask (0.00s)
=== RUN   TestRejectedCyclicReplacementKeepsExistingTask
    dag_replacement_test.go:32: rejected update removed the existing worker: task not found
--- FAIL: TestRejectedCyclicReplacementKeepsExistingTask (0.00s)
=== RUN   TestRejectedCyclicReplacementKeepsExistingTask
    dag_replacement_test.go:32: rejected update removed the existing worker: task not found
--- FAIL: TestRejectedCyclicReplacementKeepsExistingTask (0.00s)
=== RUN   TestRejectedCyclicReplacementKeepsExistingTask
    dag_replacement_test.go:32: rejected update removed the existing worker: task not found
--- FAIL: TestRejectedCyclicReplacementKeepsExistingTask (0.00s)
=== RUN   TestRejectedCyclicReplacementKeepsExistingTask
    dag_replacement_test.go:32: rejected update removed the existing worker: task not found
--- FAIL: TestRejectedCyclicReplacementKeepsExistingTask (0.00s)
FAIL
FAIL	github.com/scheduler/pkg/task	0.893s
FAIL
```
