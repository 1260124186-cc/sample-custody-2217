# 修复前故障复现（Docker）

## 项目与标准命令

项目使用 Go 1.26.2，标准构建命令为 `go build ./...`，测试命令为
`go test ./... -count=20`。

## 环境构建与编译

在当前平台使用 Docker 构建评测镜像，容器内 `go version` 为 `go1.26.2 linux/arm64`，
`go build ./...` 通过；执行 `go test ./... -count=20` 时详情数据隔离相关用例稳定失败。

## 故障触发步骤

登记并交接一个样本，读取样本详情后修改返回的样本信息、交接记录和审计事件，再读取一次详情，然后执行：

```bash
go test ./... -count=20
```

## 实际错误输出

```text
--- FAIL: TestDetailDoesNotExposeMutableHistory (0.00s)
    detail_ownership_test.go:47: detail exposed mutable state: {Sample:{ID:own-a Code:SC-OWN-001 ...} Transfers:[{... To:MUTATED ...}] Events:[{... Detail:MUTATED ...} ...]}
FAIL
FAIL	example.com/sample-custody/internal/service	1.878s
```

## 期望行为

详情接口返回的样本、交接记录和审计事件应只读；调用方修改返回值不能改变后续查询结果或历史数据。
