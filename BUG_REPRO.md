# 修复前故障复现（Docker）

## 项目与标准命令

项目使用 Go 1.26.2，标准构建命令为 `go build ./...`，测试命令为
`go test ./... -count=20`。

## 环境构建与编译

在当前平台使用 Docker 构建评测镜像，容器内 `go version` 为 `go1.26.2 linux/arm64`，
`go build ./...` 通过；执行 `go test ./... -count=20` 时检验批次封存相关用例稳定失败。

## 故障触发步骤

先把样本移入检验批次，再对包含 `in_review` 样本的批次执行完成与封存，然后执行：

```bash
go test ./... -count=20
```

## 实际错误输出

```text
--- FAIL: TestHTTPWorkflow (0.00s)
    server_test.go:29: completion response: {StatusCode:409 Body:{"error":"conflict","message":"batch completion rejected: 2 blocking, 0 warning (review-status-not-finalized,sample-cannot-seal): review-status-not-finalized: specimen must leave review before batch completion; sample-cannot-seal: sample \"http-sample-001\" cannot be sealed from status \"in_review\""}
        }
--- FAIL: TestCompletingReviewBatchSealsEverySample (0.00s)
    batch_completion_test.go:27: complete review batch: batch completion rejected: 2 blocking, 0 warning (review-status-not-finalized,sample-cannot-seal): review-status-not-finalized: specimen must leave review before batch completion; sample-cannot-seal: sample "completion-sample" cannot be sealed from status "in_review"
--- FAIL: TestCustodyWorkflow (0.00s)
    workflow_test.go:76: complete batch: batch completion rejected: 2 blocking, 0 warning (review-status-not-finalized,sample-cannot-seal): review-status-not-finalized: specimen must leave review before batch completion; sample-cannot-seal: sample "sample-test-001" cannot be sealed from status "in_review"
FAIL
FAIL	example.com/sample-custody/internal/httpapi	0.008s
FAIL	example.com/sample-custody/internal/service	0.006s
```

## 期望行为

已经进入检验批次的样本可以在批次完成时一并封存，批次完成接口应返回成功，样本状态变为已封存，不再被 `in_review` 状态阻断。
