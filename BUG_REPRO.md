# 修复前故障复现（Docker）

## 项目与标准命令

项目使用 Go 1.26.2，标准构建命令为 `go build ./...`，测试命令为
`go test ./... -count=20`。

## 环境构建与编译

在当前平台使用 Docker 构建评测镜像，容器内 `go version` 为 `go1.26.2 linux/arm64`，
`go build ./...` 通过；执行 `go test ./... -count=20` 时检验组交接相关用例稳定失败。

## 故障触发步骤

先将样本移入检验批次，再让当前保管人在检验组之间提交一次交接，然后执行：

```bash
go test ./... -count=20
```

## 实际错误输出

```text
--- FAIL: TestReviewSampleCanChangeCustodian (0.00s)
    review_transfer_test.go:24: transfer review sample: sample "review-handoff" is not available for handoff
FAIL
FAIL	example.com/sample-custody/internal/service	0.011s
```

## 期望行为

进入检验批次的样本仍应允许在检验组之间交接，交接接口应成功生成交接记录并更新当前保管人，而不是直接拒绝。
