# 修复前故障复现（Docker）

## 项目与标准命令

项目使用 Go 1.26.2，标准构建命令为 `go build ./...`，测试命令为
`go test ./... -count=20`。

## 环境构建与编译

在当前平台使用 Docker 构建评测镜像。镜像内 `go version`、`go build ./...`
和基础 `go test ./...` 均可执行；bug 环境的专项并发测试会失败。

## 故障触发步骤

在样本已经登记后，同时提交两个从同一当前保管人发往不同下一保管人的交接请求，
然后执行：

```bash
go test ./... -count=20
```

## 实际错误输出

```text
--- FAIL: TestConcurrentHandoffAcceptsOnlyOneTransfer (0.00s)
    concurrent_transfer_test.go:86: expected exactly one accepted handoff, got 2
FAIL
FAIL	example.com/sample-custody/internal/service	0.00s
```

## 期望行为

同一时刻针对同一份样本只能接受一次有效交接，另一个请求应被拒绝，最终保管人和交接记录保持一致。
