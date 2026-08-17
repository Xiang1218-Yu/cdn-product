# 分布式任务调度平台

基于Go语言实现的分布式任务调度平台，支持任务依赖DAG、执行器自动注册、故障转移和监控告警。

## 架构设计

```
┌─────────────────┐      ┌──────────────────┐
│   管理后台      │      │   调度引擎       │
│  (REST API)     │      │  (分片+故障转移) │
└─────────────────┘      └──────────────────┘
         │                        │
         │                        │
    ┌────▼────────────────────────▼────┐
    │      MongoDB (任务存储)          │
    └───────────────────────────────────┘
         │                        │
         │                        │
    ┌────▼────────────────────────▼────┐
    │      etcd (分布式锁+服务发现)    │
    └───────────────────────────────────┘
         │                        │
         │                        │
┌────────▼─────────┐      ┌──────▼──────────┐
│   执行器集群     │      │   监控告警      │
│  (gRPC通信)      │      │  (Prometheus)   │
└──────────────────┘      └─────────────────┘
```

## 核心模块

### 1. 任务管理
- CRON表达式解析
- DAG依赖关系管理
- 任务状态追踪

### 2. 调度引擎
- 任务分片调度
- 故障转移机制
- 执行器负载均衡

### 3. 执行器集群
- 自动注册与发现
- 心跳检测
- 任务执行

### 4. 监控告警
- Prometheus指标采集
- 多渠道告警通知
- 实时监控

## 技术栈

- **语言**: Go 1.21+
- **数据库**: MongoDB 7.0
- **分布式协调**: etcd v3.5+
- **RPC框架**: gRPC
- **监控**: Prometheus + Grafana
- **容器化**: Docker + Docker Compose

## 快速开始

### 1. 环境准备

```bash
# 安装依赖
go mod download

# 启动基础设施
cd deployments/docker
docker-compose up -d etcd mongodb
```

### 2. 启动调度器

```bash
go run cmd/scheduler/main.go --config configs/config.yaml
```

### 3. 启动执行器

```bash
go run cmd/executor/main.go --id executor-1 --capacity 10
```

### 4. 启动管理后台

```bash
go run cmd/admin/main.go --config configs/config.yaml
```

### 5. 使用Docker Compose一键部署

```bash
cd deployments/docker
docker-compose up -d
```

## API接口

### 创建任务

```bash
POST /api/v1/tasks
Content-Type: application/json

{
  "name": "daily-report",
  "command": "python /scripts/report.py",
  "cron_expr": "0 0 9 * * *",
  "params": {
    "type": "daily"
  },
  "timeout": 3600,
  "max_retries": 3
}
```

### 查询任务列表

```bash
GET /api/v1/tasks/list
```

### 查询任务详情

```bash
GET /api/v1/tasks/get?id=task-id
```

### 删除任务

```bash
DELETE /api/v1/tasks/delete?id=task-id
```

### 查询任务日志

```bash
GET /api/v1/logs?task_id=task-id
```

## 执行器SDK使用

```go
package main

import (
    "context"
    "github.com/scheduler/pkg/sdk"
    "go.uber.org/zap"
)

func main() {
    logger, _ := zap.NewProduction()
    
    // 创建执行器
    executor := sdk.NewSimpleExecutor(
        "executor-1",
        "localhost:9090",
        10,
        func(ctx context.Context, command string, params map[string]string) (string, error) {
            // 执行任务逻辑
            return "task executed", nil
        },
        logger,
    )
    
    // 注册到调度器
    if err := executor.Register("localhost:9090"); err != nil {
        logger.Fatal("failed to register", zap.Error(err))
    }
}
```

## 监控指标

访问 `http://localhost:2112/metrics` 查看Prometheus指标：

- `scheduler_tasks_total`: 任务总数
- `scheduler_tasks_success_total`: 成功任务数
- `scheduler_tasks_failed_total`: 失败任务数
- `scheduler_tasks_duration_seconds`: 任务执行时长
- `scheduler_executors_active`: 活跃执行器数
- `scheduler_executors_load`: 执行器负载

## 配置说明

编辑 `configs/config.yaml` 文件：

```yaml
server:
  scheduler:
    host: "0.0.0.0"
    port: 8080      # HTTP API端口
    grpc_port: 9090 # gRPC端口

etcd:
  endpoints:
    - "http://localhost:2379"
  dial_timeout: 5s

mongodb:
  uri: "mongodb://localhost:27017"
  database: "scheduler"
  timeout: 10s

retry:
  max_attempts: 3
  initial_delay: 1s
  max_delay: 30s
```

## 目录结构

```
scheduler/
├── cmd/                    # 主程序入口
│   ├── scheduler/          # 调度器
│   ├── executor/           # 执行器
│   └── admin/              # 管理后台
├── pkg/                    # 公共包
│   ├── task/               # 任务管理
│   ├── scheduler/          # 调度引擎
│   ├── executor/           # 执行器服务
│   ├── retry/              # 重试机制
│   ├── log/                # 日志模块
│   ├── monitor/            # 监控告警
│   ├── sdk/                # 执行器SDK
│   └── api/                # REST API
├── internal/               # 内部包
│   ├── config/             # 配置管理
│   ├── store/              # 数据存储
│   ├── lock/               # 分布式锁
│   └── discovery/          # 服务发现
├── api/proto/              # gRPC协议定义
├── configs/                # 配置文件
└── deployments/docker/     # Docker部署文件
```

## 特性

✅ 支持CRON表达式定时任务  
✅ 支持任务依赖DAG  
✅ 执行器自动注册与发现  
✅ 任务分片调度  
✅ 故障自动转移  
✅ 失败重试机制  
✅ 任务执行日志追踪  
✅ Prometheus监控集成  
✅ 多渠道告警通知  
✅ Docker容器化部署  

## 开发计划

- [ ] Web管理界面
- [ ] 任务执行历史统计
- [ ] 任务依赖可视化
- [ ] 执行器资源监控
- [ ] 任务优先级队列
- [ ] 任务执行超时控制

## 许可证

MIT License