# 修复前故障复现（Docker）

## 项目与标准命令

项目使用 Go 1.26.2，标准构建命令为 `go build ./...`，测试命令为
`go test ./... -count=20`。

## 环境构建与编译

在当前平台使用 Docker 构建评测镜像，容器内 `go version` 为 `go1.26.2 linux/arm64`，
`go build ./...` 通过；执行 `go test ./... -count=20` 时按当前保管人检索相关用例稳定失败。

## 故障触发步骤

把一份样本转交到检验组后，用新的当前保管人执行样本检索，然后运行：

```bash
go test ./... -count=20
```

## 实际错误输出

```text
--- FAIL: TestSearchFindsReviewSampleByCurrentHolder (0.00s)
    review_search_test.go:23: unexpected search results: []
FAIL
FAIL	example.com/sample-custody/internal/service	0.002s
```

## 期望行为

样本转交到检验组后，按新的当前保管人检索应能查到该样本，列表与接口返回的保管人信息都应是检验组。
