# 修复前故障复现（Docker）

## 项目与标准命令

项目使用 Go 1.26.2，标准构建命令为 `go build ./...`，测试命令为
`go test ./... -count=20`。

## 环境构建与编译

在当前平台使用 Docker 构建评测镜像，容器内 `go version` 为 `go1.26.2 linux/arm64`，
`go build ./...` 通过；执行 `go test ./... -count=20` 时开放批次清单相关用例稳定崩溃。

## 故障触发步骤

创建一个仍未完成的检验批次，直接请求该批次的清单，然后执行：

```bash
go test ./... -count=20
```

## 实际错误输出

```text
--- FAIL: TestOpenBatchManifestDoesNotPanic (0.00s)
panic: runtime error: invalid memory address or nil pointer dereference
...
example.com/sample-custody/internal/service.(*QueryService).Manifest
	/app/internal/service/query_service.go:43 +0xc0
```

## 期望行为

未完成批次也能正常生成清单，返回开放状态文本，不应因缺失完成时间而崩溃。
