# 修复前故障复现（Docker）

## 项目与标准命令

项目使用 Go 1.26.2，标准构建命令为 `go build ./...`，测试命令为
`go test ./... -count=20`。

## 环境构建与编译

在当前平台使用 Docker 构建评测镜像，容器内 `go version` 为 `go1.26.2 linux/arm64`，
`go build ./...` 通过；执行 `go test ./... -count=20` 时批次数据隔离相关用例稳定失败。

## 故障触发步骤

创建包含多个样本的检验批次，修改返回批次里的样本编号列表后再次查询同一批次，然后执行：

```bash
go test ./... -count=20
```

## 实际错误输出

```text
--- FAIL: TestBatchOwnsItsSampleIDSlice (0.00s)
    batch_ownership_test.go:51: mutating returned batch changed stored sample ids: [mutated-id slice-b]
FAIL
FAIL	example.com/sample-custody/internal/service	1.419s
```

## 期望行为

调用方修改返回批次中的样本编号列表，不应影响存储中的批次数据；再次查询应始终返回原始编号。
