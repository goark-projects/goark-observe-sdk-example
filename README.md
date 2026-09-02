# Goark Observe SDK Example

[简体中文](README.zh-CN.md)

This project is a runnable reference for `goark.dev/observe-sdk`. It creates a provider, emits correlated traces, metrics, logs, and events, flushes pending telemetry, and shuts the provider down safely.

## Requirements

- Go 1.25 or later

## Run

```bash
go run ./cmd/sdk-example
```

The command writes one JSON object per exported signal. The console exporter exists only to make the SDK behavior visible; production applications should use a dedicated exporter such as `goark.dev/observe-exporter-otel`, `goark.dev/observe-exporter-prometheus`, or a private implementation.

## What It Demonstrates

- Resource and instrumentation-scope configuration.
- Trace creation, status, attributes, and span events.
- Logs and events correlated with the active trace context.
- A bounded-label counter and explicit-bound latency histogram.
- Explicit `ForceFlush` and idempotent `Shutdown` lifecycle handling.
- A custom exporter implementing all four signal interfaces.

## Layout

```text
cmd/sdk-example/             executable entry point
internal/example/            signal creation and lifecycle orchestration
internal/consoleexporter/    synchronous JSON-lines exporter for demonstration
```

## Validation

```bash
go test ./...
go test -race ./...
go vet ./...
```

## Production Notes

- Keep metric attributes bounded. Do not use request IDs, user IDs, raw URLs, SQL text, or error messages as metric labels.
- Exporters on synchronous paths must be non-blocking. Use a bounded processor when an exporter can block.
- Always flush and shut down the provider during application termination.
- Never place credentials, tokens, personal data, or unrestricted request payloads in telemetry.

Licensed under Apache License 2.0.
