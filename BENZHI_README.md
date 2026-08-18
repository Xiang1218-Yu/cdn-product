# CDN Product 分布式任务调度平台

本项目提供一个 Go 实现的分布式任务调度平台：负责创建和管理任务、按计划调度任务、将任务派发给执行器集群，并提供任务日志、监控与管理接口。调度器、执行器和管理后台分别位于 `cmd/scheduler`、`cmd/executor` 和 `cmd/admin`。

## 环境要求

- Go 1.22 或更高版本（项目模块声明为 Go 1.21，可由 Go 1.22 构建）
- Docker（仅在使用镜像构建或启动依赖服务时需要）

## 构建

在项目根目录执行：

```bash
go build ./...
```

首次准备依赖时可执行：

```bash
go mod download
```

## 运行

调度器和管理后台依赖 Etcd、MongoDB 等基础服务。可先启动本地依赖：

```bash
cd deployments/docker
docker-compose up -d etcd mongodb
cd ../..
```

分别启动各组件：

```bash
# 启动调度器
go run cmd/scheduler/main.go --config configs/config.yaml

# 启动一个执行器
go run cmd/executor/main.go --id executor-1 --capacity 10

# 启动管理后台
go run cmd/admin/main.go --config configs/config.yaml
```

如需使用 Docker Compose 启动整套服务：

```bash
cd deployments/docker
docker-compose up -d
```

## 测试

运行全部 Go 测试：

```bash
go test ./...
```

## Benzhi 评测镜像

使用项目根目录的辅助脚本构建评测镜像：

```bash
./build_benzhi_docker.sh cdn-product
```

可选的第二个参数用于指定镜像平台，例如：

```bash
./build_benzhi_docker.sh cdn-product linux/amd64
```

构建完成后进入容器：

```bash
docker run -it cdn-product:latest
```
