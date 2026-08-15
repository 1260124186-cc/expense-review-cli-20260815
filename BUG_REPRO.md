# 修复前故障复现（Docker）

## 文档范围
本文件仅记录修复前 base 状态的公开可观察故障，用于在 Docker 环境中复现。
不描述修复后的行为、根因或实现方案。

## 项目与标准命令
Expense Review CLI 用于根据分类限额和收据规则审核本地 JSON 报销批次。在仓库根目录可执行以下标准命令：

```sh
go build ./...
go run ./cmd/expense-review --input examples/claims.json
go test ./...
```

本故障通过 `go test ./...` 稳定触发。`go run ./cmd/expense-review --input examples/claims.json` 在修复前会输出对示例报销的错误拒绝结果。

## 环境构建与编译
以下命令已经实际执行。`linux/amd64` 和 `linux/arm64` 的镜像构建以及容器内 `go build ./...` 均成功：

```sh
docker build --platform linux/amd64 -f benzhi.Dockerfile -t expense-review-cli-bug003-base:amd64 .
docker run --rm --platform linux/amd64 expense-review-cli-bug003-base:amd64 go build ./...
docker build --platform linux/arm64 -f benzhi.Dockerfile -t expense-review-cli-bug003-base:arm64 .
docker run --rm --platform linux/arm64 expense-review-cli-bug003-base:arm64 go build ./...
```

目标故障在下节的测试命令中触发，不是镜像构建或容器内编译失败。

## 故障触发步骤
在仓库根目录执行以下命令：

```sh
docker run --rm --platform linux/amd64 expense-review-cli-bug003-base:amd64 go test ./...
```

## 实际错误输出
上述命令在修复前 base 状态的完整输出如下：

```text
?   	github.com/1260124186-cc/expense-review-cli-20260815/cmd/expense-review	[no test files]
--- FAIL: TestDefaultPolicyProvidesCategoryCaps (0.00s)
    policy_test.go:7: meal cap = 0, want 7500
FAIL
FAIL	github.com/1260124186-cc/expense-review-cli-20260815/internal/domain	0.050s
--- FAIL: TestReviewAndRender (0.03s)
    review_test.go:34: rendered output missing "meal-1=review (receipt required)":
        period=2026-08 total_cents=87000
        meal-1=rejected (category cap exceeded)
        trip-1=rejected (category cap exceeded)
--- FAIL: TestReviewUsesDefaultCapsWhenPolicyOmitsCategoryCaps (0.00s)
    review_test.go:69: rendered output = "period=2026-08 total_cents=4200\nmeal-1=rejected (category cap exceeded)\n", want an approved meal claim
FAIL
FAIL	github.com/1260124186-cc/expense-review-cli-20260815/internal/service	0.120s
?   	github.com/1260124186-cc/expense-review-cli-20260815/internal/store	[no test files]
FAIL
```

该命令退出状态为 1。

## 期望行为
在未提供完整分类限额的输入下，报销应继续按标准规则审核：满足分类限额且未达到收据要求的记录应通过，达到收据要求的记录应进入人工复核，而不是因分类限额错误被拒绝。
