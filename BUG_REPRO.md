# 修复前故障复现（Docker）

## 文档范围

本文件仅记录修复前 base 状态的公开可观察故障，用于在 Docker 环境中复现。
不描述修复后的行为、根因或实现方案。

## 项目与标准命令

Expense Review CLI 用于根据费用类别和票据规则审核本地 JSON 报销批次。在仓库根目录可执行：

```sh
go build ./...
go run ./cmd/expense-review --input examples/claims.json
go test ./...
```

其中 `go test ./...` 在修复前状态会稳定失败；运行命令不带 `--output` 参数时用于输出审核结果到标准输出。

## 环境构建与编译

已实际执行以下 linux/amd64 命令，并确认镜像构建和容器内编译成功：

```sh
docker build --platform linux/amd64 -f benzhi.Dockerfile -t expense-review-bug005-base-amd64 .
docker run --rm --platform linux/amd64 expense-review-bug005-base-amd64 go build ./...
```

已实际执行以下 linux/arm64 命令，并确认镜像构建和容器内编译成功：

```sh
docker build --platform linux/arm64 -f benzhi.Dockerfile -t expense-review-bug005-base-arm64 .
docker run --rm --platform linux/arm64 expense-review-bug005-base-arm64 go build ./...
```

两个平台均已实际启动样例 CLI：

```sh
docker run --rm --platform linux/amd64 expense-review-bug005-base-amd64
docker run --rm --platform linux/arm64 expense-review-bug005-base-arm64
```

## 故障触发步骤

在仓库根目录先按上节命令构建 linux/arm64 镜像，然后执行：

```sh
docker run --rm --platform linux/arm64 expense-review-bug005-base-arm64 go test ./...
```

该命令已连续执行 20 次，均稳定失败。

## 实际错误输出

```text
?   	github.com/1260124186-cc/expense-review-cli-20260815/cmd/expense-review	[no test files]
?   	github.com/1260124186-cc/expense-review-cli-20260815/internal/domain	[no test files]
--- FAIL: TestWritePublishesRenderedReview (0.00s)
    review_test.go:50: published output = "", want "period=2026-08 total_cents=1\n"
FAIL
FAIL	github.com/1260124186-cc/expense-review-cli-20260815/internal/service	0.003s
?   	github.com/1260124186-cc/expense-review-cli-20260815/internal/store	[no test files]
FAIL
```

## 期望行为

当命令使用 `--output` 发布审核结果并返回成功时，目标文件应包含完整的审核文本，不应为空文件；如果写入或发布失败，命令应返回失败结果。
