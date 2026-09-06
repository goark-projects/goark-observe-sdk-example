# Goark Observe SDK 示例

[English](README.md)

本项目是 `goark.dev/observe-sdk` 的可运行参考实现，展示 Provider 创建、trace、metric、log、event 关联、待处理数据刷新以及 Provider 安全关闭。

## 环境要求

- Go 1.26 或更高版本
- 示例唯一使用字节跳动 Sonic 处理 JSON。

## 运行

```bash
go run ./cmd/sdk-example
```

命令会为每条导出的观测信号输出一个 JSON 对象。控制台 exporter 仅用于直观展示 SDK 行为；生产应用应使用 `goark.dev/observe-exporter-otel`、`goark.dev/observe-exporter-prometheus` 或自研 exporter。

## 示例内容

- Resource 和 instrumentation scope 配置。
- Trace 创建、状态、属性和 span event。
- Log、event 与当前 trace 上下文关联。
- 使用有界标签的 counter 和显式边界延迟 histogram。
- 显式 `ForceFlush` 和幂等 `Shutdown` 生命周期处理。
- 同时实现四类信号接口的自定义 exporter。

## 目录结构

```text
cmd/sdk-example/             可执行程序入口
internal/example/            信号生成和生命周期编排
internal/consoleexporter/    示例使用的同步 JSON 行 exporter
```

## 验证

```bash
go test ./...
go test -race ./...
go vet ./...
```

## 生产注意事项

- 指标属性必须保持低基数，不要把请求 ID、用户 ID、原始 URL、SQL 文本或错误消息作为指标标签。
- 同步路径中的 exporter 必须非阻塞；可能阻塞的 exporter 应放入有界 processor。
- 应用终止时必须刷新并关闭 Provider。
- 不要把凭据、令牌、个人数据或未限制的请求内容写入 telemetry。

使用 Apache License 2.0。
