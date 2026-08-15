# 修复前故障复现（Docker）

## 文档范围
本文件仅记录修复前 base 状态的公开可观察故障，用于在 Docker 环境中复现。
不描述修复后的行为、根因或实现方案。

## 项目与标准命令
Expense Review CLI 用于根据费用类别和票据规则审核 JSON 格式的员工报销批次。
在仓库根目录可以执行以下标准命令：

```sh
go build ./...
go run ./cmd/expense-review --input examples/claims.json
go test ./...
```

其中 `go test ./...` 是本故障的稳定触发命令；修复前该命令预期失败。

## 环境构建与编译
已在修复前 base 状态实际执行以下 linux/amd64 命令：

```sh
docker build --platform linux/amd64 -f benzhi.Dockerfile -t expense-review-base-bug001:verify-amd64 .
docker run --rm --platform linux/amd64 expense-review-base-bug001:verify-amd64 go build ./...
```

已在修复前 base 状态实际执行以下 linux/arm64 命令：

```sh
docker build --platform linux/arm64 -f benzhi.Dockerfile -t expense-review-base-bug001:verify-arm64 .
docker run --rm --platform linux/arm64 expense-review-base-bug001:verify-arm64 go build ./...
```

两个平台的镜像构建和容器内 `go build ./...` 均成功；目标故障通过下一节的测试命令触发。

## 故障触发步骤
在仓库根目录构建 linux/amd64 镜像后，执行：

```sh
docker run --rm --platform linux/amd64 expense-review-base-bug001:verify-amd64 go test ./...
```

该命令在修复前 base 状态稳定失败。

## 实际错误输出

```text
?   	github.com/1260124186-cc/expense-review-cli-20260815/cmd/expense-review	[no test files]
--- FAIL: TestBatchRejectsDuplicateClaimIDs (0.00s)
    batch_test.go:15: Validate() error = nil, want duplicate claim rejection
FAIL
FAIL	github.com/1260124186-cc/expense-review-cli-20260815/internal/domain	0.038s
--- FAIL: TestReviewRejectsDuplicateClaims (0.00s)
    review_test.go:65: ReviewAndRender() error = nil, want duplicate input rejection
FAIL
FAIL	github.com/1260124186-cc/expense-review-cli-20260815/internal/service	0.046s
?   	github.com/1260124186-cc/expense-review-cli-20260815/internal/store	[no test files]
FAIL
```

命令退出状态：`1`。

## 期望行为
当同一批报销记录包含重复编号时，CLI 应拒绝该输入并返回清楚错误，不生成可继续消费的审核结果；具有唯一编号的正常批次仍应正常审核。
