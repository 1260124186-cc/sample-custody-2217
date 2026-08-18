# 修复前故障复现（Docker）

## 项目与标准命令

项目使用 Go 1.26.2，标准构建命令为 `go build ./...`，测试命令为
`go test ./... -count=20`。

## 环境构建与编译

在当前平台使用 Docker 构建评测镜像，容器内 `go version` 为 `go1.26.2 linux/arm64`，
`go build ./...` 通过；执行 `go test ./... -count=20` 时编码唯一性相关用例稳定失败。

## 故障触发步骤

先登记编码为 `SC-CASE-001` 的样本，再登记编码为 `sc-case-001` 的样本，然后执行：

```bash
go test ./... -count=20
```

## 实际错误输出

```text
--- FAIL: TestSampleCodeUniquenessIsCaseInsensitive (0.00s)
    code_uniqueness_test.go:32: expected conflict for case-insensitive duplicate code, got <nil>
FAIL
FAIL	example.com/sample-custody/internal/service	1.300s
```

## 期望行为

样本编码唯一性应不区分大小写，大小写不同的重复编码应被拒绝，登记、检索和批次校验保持同一套判断口径。
