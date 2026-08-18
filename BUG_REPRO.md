# 修复前故障复现（Docker）

## 项目与标准命令

项目使用 Go 1.26.2，标准构建命令为 `go build ./...`，测试命令为
`go test ./... -count=20`。

## 环境构建与编译

在当前平台使用 Docker 构建评测镜像，容器内 `go version` 为 `go1.26.2 linux/arm64`，
`go build ./...` 通过；执行 `go test ./... -count=20` 时封存汇总统计相关用例稳定失败。

## 故障触发步骤

登记并封存一份样本后，读取样本汇总统计，然后执行：

```bash
go test ./... -count=20
```

## 实际错误输出

```text
--- FAIL: TestSummaryCountsSealedSampleForItsOrigin (0.00s)
    sealed_summary_test.go:19: unexpected sealed summary: {SampleCount:1 RegisteredCount:0 InReviewCount:0 SealedCount:0 OpenBatchCount:0 CompletedBatches:0 TransferCount:0 Origins:[{Origin:archive-lab SampleCount:1 SealedCount:0 OpenCount:1}] Holders:[{Holder:intake SampleCount:1 SealedCount:1}]}
FAIL
FAIL	example.com/sample-custody/internal/service	0.003s
```

## 期望行为

样本封存后，汇总统计应正确体现已封存数量、来源归档状态和当前保管人，导出的状态文本与批次清单也应保持一致。
