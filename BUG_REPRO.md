# 第 7 题 Bug 复现说明

## Bug 是什么
调度器释放 leader 租约后没有撤销自身的 leader 资格，后续调度循环仍可能继续派发 pending 任务。与此同时，etcd 租约释放仅解锁而未关闭 session，可能遗留租约相关资源。

## 如何触发
建立一个处于 leader 状态的调度器，释放其 leader 租约后检查调度资格是否已经撤销。

## 运行指令
```bash
go test ./pkg/scheduler -run '^TestReleasedLeaderLeaseStopsScheduling$' -count=1
```

## 错误信息
失败输出说明租约已经释放，但调度器仍保留 leader 标记，仍具备继续派发任务的资格。

## 错误堆栈
```text
--- FAIL: TestReleasedLeaderLeaseStopsScheduling (0.00s)
    engine_leadership_test.go:31: released leader lease still leaves scheduler eligible to dispatch tasks
FAIL
FAIL	github.com/scheduler/pkg/scheduler	0.331s
FAIL
```
