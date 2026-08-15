# Docker 交付说明

本项目提供可独立构建、编译和运行的 Docker 环境，用于复验
Expense Review CLI 的源码状态。镜像保留完整 Go 工具链，所有校验均在
容器内针对当前源码执行；不依赖宿主机编译出的二进制，也不需要数据库、
消息队列或其他应用服务。

## 交付文件与环境约束

| 文件 | 作用 |
| --- | --- |
| `benzhi.Dockerfile` | 固定 Go 构建环境，下载依赖并在镜像构建阶段执行 `go build ./...`。 |
| `build_benzhi_docker.sh` | 对指定平台构建镜像，并在容器内完成编译和实际启动校验。 |
| `.dockerignore` | 排除 Git 元数据、补丁、缓存、日志和本地编译产物，确保构建上下文只包含项目交付所需内容。 |
| `examples/claims.json` | 默认启动命令使用的本地示例输入。 |

Dockerfile 使用 `golang:1.26.2-bookworm`，并设置
`GOTOOLCHAIN=local`。Go 工具链版本与 `go.mod` 中的 Go 语言版本共同
构成可复现环境的一部分；请勿因宿主机版本不同而修改 Dockerfile、
`go.mod` 或工具链设置。

该镜像是用于源码验收的构建环境，而非只保留可执行文件的精简运行镜像。
保留 Go 工具链是必要条件，便于在 `linux/amd64` 和 `linux/arm64` 容器中
重新编译当前源码。

## 前置条件

在仓库根目录执行以下命令：

```sh
cd /path/to/expense-review-cli-20260815
```

本机需要可用的 Docker 服务，并支持 `docker build` 和 `docker run`。
宿主机无需安装 Go。若本机 CPU 架构与目标平台不一致，Docker 运行时必须
已提供对应的平台模拟能力；不要通过在 Dockerfile 中写死 CPU 架构来规避
该问题。

## 标准双架构验收

分别执行两个目标平台的完整校验：

```sh
./build_benzhi_docker.sh linux/amd64
./build_benzhi_docker.sh linux/arm64
```

脚本未传参数时默认验证 `linux/amd64`：

```sh
./build_benzhi_docker.sh
```

每次执行会按以下顺序完成：

1. 使用 `benzhi.Dockerfile` 构建 `expense-review-cli:benzhi` 镜像。
2. 在新容器内执行 `go build ./...`，确认容器中的完整源码可以编译。
3. 在另一新容器内执行镜像默认命令，使用 `examples/claims.json` 实际启动 CLI。

两个平台上的命令均以退出码 `0` 结束，才表示 Docker 环境验收通过。构建日志
中出现 `go mod download` 和 `go build ./...` 是预期行为。

## 手工复验

需要定位环境问题时，可将平台替换为待验证的值并逐步执行：

```sh
docker build --platform linux/amd64 \
  -f benzhi.Dockerfile \
  -t expense-review-cli:benzhi .

docker run --rm --platform linux/amd64 \
  expense-review-cli:benzhi go build ./...

docker run --rm --platform linux/amd64 \
  expense-review-cli:benzhi
```

默认启动命令等价于：

```sh
go run ./cmd/expense-review --input examples/claims.json
```

项目的公开回归测试也可在同一镜像中执行：

```sh
docker run --rm --platform linux/amd64 \
  expense-review-cli:benzhi go test ./...
```

## 常见错误处理

| 现象 | 处理方式 |
| --- | --- |
| `Cannot connect to the Docker daemon` | 启动 Docker 服务后，在仓库根目录重新执行原命令。 |
| 非本机平台无法构建或容器无法启动 | 检查 Docker 的多架构模拟支持，或在对应 CPU 架构的环境中复验；不要删除 `--platform`，也不要在 Dockerfile 中固定单一架构。 |
| `go mod download` 失败 | 检查 Docker 构建环境到 Go 模块源的网络连通性后重试。不要将宿主机模块缓存或下载后的依赖目录复制进镜像。 |
| Go 版本与预期不一致 | 以 Dockerfile 固定的工具链、`go.mod` 和 `GOTOOLCHAIN=local` 为准；不要用宿主机 Go 版本替代容器校验。 |
| 本地文件意外进入构建上下文 | 检查 `.dockerignore`。Git 元数据、补丁、日志、缓存和编译产物不得进入镜像或替代源码构建。 |

交付复验应始终使用上述容器内命令和当前仓库源码。不要通过复制宿主机二进制、
跳过容器内编译、替换为未记录的 Go 版本，或修改验证命令来制造通过结果。
