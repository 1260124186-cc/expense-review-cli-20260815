# Expense Review CLI

Expense Review CLI 读取员工报销单 JSON，按类别额度和票据规则输出审核结果。
项目本地运行，不依赖数据库或外部服务。

## 标准命令

在仓库根目录执行。项目使用 Go 1.26.2。

```sh
go build ./...                                      # 编译
go run ./cmd/expense-review --input examples/claims.json  # 启动
go test ./...                                       # 测试
```

如需将审核结果写入文件：

```sh
go run ./cmd/expense-review \
  --input examples/claims.json \
  --output review.txt
```

## Docker 验收

`benzhi.Dockerfile` 使用固定 Go 工具链并设置 `GOTOOLCHAIN=local`。
镜像保留 Go 工具链，确保验证在容器内基于当前源码完成。

分别验证两个交付平台：

```sh
./build_benzhi_docker.sh linux/amd64
./build_benzhi_docker.sh linux/arm64
```

每次执行都会构建镜像、在容器内运行 `go build ./...`，再使用
`examples/claims.json` 启动 CLI。两个命令均以退出码 `0` 结束即表示通过。

也可以手工验证：

```sh
docker build --platform linux/amd64 \
  -f benzhi.Dockerfile \
  -t expense-review-cli:benzhi .

docker run --rm --platform linux/amd64 \
  expense-review-cli:benzhi go build ./...

docker run --rm --platform linux/amd64 \
  expense-review-cli:benzhi go test ./...

docker run --rm --platform linux/amd64 \
  expense-review-cli:benzhi
```

将命令中的 `linux/amd64` 替换为 `linux/arm64` 可进行另一平台的手工复验。
