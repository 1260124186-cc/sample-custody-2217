# 修复前故障复现（Docker）

## 项目与标准命令

项目使用 Go 1.26.2，标准构建命令为 `go build ./...`，测试命令为
`go test ./... -count=20`。

## 环境构建与编译

在当前平台使用 Docker 构建评测镜像，容器内 `go version` 为 `go1.26.2 linux/arm64`，
`go build ./...` 通过；执行 `go test ./... -count=20` 时错误分类相关用例稳定失败。

## 故障触发步骤

先登记一个样本，再使用相同样本编码提交第二次登记，然后执行：

```bash
go test ./... -count=20
```

## 实际错误输出

```text
--- FAIL: TestDomainErrorsStayClassified (0.00s)
    error_contract_test.go:28: expected conflict domain error, got create sample failed: sample code "SC-DOM-001" already exists
FAIL
FAIL	example.com/sample-custody/internal/service	0.935s
```

## 期望行为

重复样本编码应返回可识别的业务冲突错误，客户端能明确区分冲突与普通系统失败。
