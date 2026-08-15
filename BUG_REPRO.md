# 修复前故障复现（Docker）

## 文档范围
本文件仅记录修复前 base 状态的公开可观察故障，用于在 Docker 环境中复现。
不描述修复后的行为、根因或实现方案。

## 项目与标准命令
Expense Review CLI 用于根据分类限额和凭据规则审核本地 JSON 报销批次。在仓库根目录执行：

```sh
go build ./...
go run ./cmd/expense-review --input examples/claims.json
go test ./...
```

修复前，`go build ./...` 和示例运行可以完成；`go test ./...` 会稳定暴露审核后原始报销单顺序被改变的故障。

## 环境构建与编译
修复前 base 状态在两个平台均已实际完成镜像构建和容器内编译：

```sh
docker build --platform linux/amd64 -f benzhi.Dockerfile -t expense-review-cli-20260815-bug-002-base:amd64 .
docker run --rm --platform linux/amd64 expense-review-cli-20260815-bug-002-base:amd64 go build ./...
docker build --platform linux/arm64 -f benzhi.Dockerfile -t expense-review-cli-20260815-bug-002-base:arm64 .
docker run --rm --platform linux/arm64 expense-review-cli-20260815-bug-002-base:arm64 go build ./...
```

两个平台的镜像构建和容器内 `go build ./...` 均成功；目标故障在下节的测试命令中触发。

## 故障触发步骤
在仓库根目录先按上节构建 linux/amd64 镜像，再执行：

```sh
docker run --rm --platform linux/amd64 expense-review-cli-20260815-bug-002-base:amd64 go test ./...
```

该命令已连续执行 20 次，均稳定失败。

## 实际错误输出

```text
?   	github.com/1260124186-cc/expense-review-cli-20260815/cmd/expense-review	[no test files]
--- FAIL: TestReviewBatchDoesNotReorderClaims (0.00s)
    review_test.go:18: claim order changed to "a-first", want "z-last"
FAIL
FAIL	github.com/1260124186-cc/expense-review-cli-20260815/internal/domain	0.031s
--- FAIL: TestReviewDoesNotReorderRepositoryClaims (0.00s)
    review_test.go:70: repository claim order changed to "a-first", want "z-last"
FAIL
FAIL	github.com/1260124186-cc/expense-review-cli-20260815/internal/service	0.039s
?   	github.com/1260124186-cc/expense-review-cli-20260815/internal/store	[no test files]
FAIL

退出状态: 1
```

## 期望行为
同一批报销记录完成审核后，调用方和已保存批次仍应保持提交时的原始顺序；审核结果可以按既定规则输出，且正常示例输入应继续返回完整的审核文本和成功状态。
