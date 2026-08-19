# 第 6 题 Bug 复现说明

## Bug 是什么
执行器收到任务请求中的短超时配置后，没有把该时限传入实际执行上下文。任务只能等到更长的父请求截止时间结束，并且最终被记录为普通失败，导致调度端无法及时依据超时终态回收执行容量。

## 如何触发
构造一个父上下文时限为 80ms、任务自身时限为 8ms 的执行请求，并让任务函数等待执行上下文结束。

## 运行指令
```bash
go test ./pkg/executor -run '^TestExecutorHonorsRequestDeadlineAndRecordsTimeout$' -count=1
```

## 错误信息
失败输出表明执行持续到了父请求的约 80ms，而不是任务请求声明的 8ms；因此短时限没有在执行端生效。

## 错误堆栈
```text
--- FAIL: TestExecutorHonorsRequestDeadlineAndRecordsTimeout (0.08s)
    task_runner_deadline_test.go:39: request timeout was ignored: execution lasted 81.147375ms
FAIL
FAIL	github.com/scheduler/pkg/executor	0.691s
FAIL
```
