# 修复前故障复现（Docker）

## 文档范围
本文件仅记录修复前 base 状态的公开可观察故障，用于在 Docker 环境中复现。
不描述修复后的行为、根因或实现方案。

## 项目与标准命令
Expense Review CLI 用于在本地读取 JSON 格式的报销批次，并按照分类限额和凭证规则输出审核结果。

在仓库根目录可执行以下标准命令：

```sh
go build ./...
go run ./cmd/expense-review --input examples/claims.json
go test ./...
```

前两条命令可以正常完成。`go test ./...` 会触发本文件记录的取消请求故障。

## 环境构建与编译
以下 Docker 命令已在修复前 base 状态实际执行：

```sh
docker build --platform linux/amd64 -f benzhi.Dockerfile -t expense-review-cli:bug-004-base-amd64 .
docker run --rm --platform linux/amd64 expense-review-cli:bug-004-base-amd64 go build ./...
docker build --platform linux/arm64 -f benzhi.Dockerfile -t expense-review-cli:bug-004-base-arm64 .
docker run --rm --platform linux/arm64 expense-review-cli:bug-004-base-arm64 go build ./...
```

linux/amd64 与 linux/arm64 的镜像构建和容器内 `go build ./...` 均成功。两个镜像中的示例命令也能正常启动：

```sh
docker run --rm --platform linux/amd64 expense-review-cli:bug-004-base-amd64
docker run --rm --platform linux/arm64 expense-review-cli:bug-004-base-arm64
```

## 故障触发步骤
在仓库根目录执行：

```sh
go test ./...
```

该命令会稳定失败，表明调用方已经取消审核请求时，程序仍然返回了完成结果而不是取消错误。

## 实际错误输出
修复前 base 状态一次完整的实际失败输出如下：

```text
?   	github.com/1260124186-cc/expense-review-cli-20260815/cmd/expense-review	[no test files]
?   	github.com/1260124186-cc/expense-review-cli-20260815/internal/domain	[no test files]
--- FAIL: TestReviewHonorsCanceledContext (0.00s)
    review_test.go:66: Review() error = nil, want canceled context
FAIL
FAIL	github.com/1260124186-cc/expense-review-cli-20260815/internal/service	1.391s
?   	github.com/1260124186-cc/expense-review-cli-20260815/internal/store	[no test files]
FAIL
```

## 期望行为
当调用方在审核开始前或处理中取消请求时，审核应尽快停止并返回取消错误，不应返回完整审核结果或继续产生可供下游使用的完成输出。未取消的请求应继续输出正常审核结果。
