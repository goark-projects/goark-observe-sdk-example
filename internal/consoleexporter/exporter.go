// Package consoleexporter 提供仅用于示例和本地调试的 JSON 行 exporter。
package consoleexporter

import (
	"context"
	"encoding/json"
	"io"
	"sync"

	"goark.dev/observe"
)

// Exporter 把四类观测信号以 JSON 行写入目标输出。
type Exporter struct {
	mu      sync.Mutex
	encoder *json.Encoder
}

type record struct {
	Signal       string            `json:"signal"`
	Name         string            `json:"name"`
	Resource     string            `json:"resource,omitempty"`
	Scope        string            `json:"scope,omitempty"`
	TraceID      string            `json:"trace_id,omitempty"`
	SpanID       string            `json:"span_id,omitempty"`
	ParentSpanID string            `json:"parent_span_id,omitempty"`
	Level        string            `json:"level,omitempty"`
	Message      string            `json:"message,omitempty"`
	Value        string            `json:"value,omitempty"`
	Points       int               `json:"points,omitempty"`
	Attributes   map[string]string `json:"attributes,omitempty"`
}

// New 创建并发安全的控制台 exporter；writer 为空时丢弃输出。
func New(writer io.Writer) *Exporter {
	if writer == nil {
		writer = io.Discard
	}
	return &Exporter{encoder: json.NewEncoder(writer)}
}

// Descriptor 返回示例 exporter 的稳定能力描述。
func (*Exporter) Descriptor() observe.ExporterDescriptor {
	return observe.ExporterDescriptor{
		Name:      "console-json",
		Version:   "1.0.0",
		Signals:   observe.SignalAll,
		Stability: observe.StabilityExperimental,
		Capabilities: observe.ExporterCapabilities{
			Push:                  true,
			CumulativeTemporality: true,
			Histogram:             true,
		},
	}
}

// ForceFlush 完成同步 exporter 的生命周期契约。
func (*Exporter) ForceFlush(context.Context) error { return nil }

// Shutdown 完成同步 exporter 的生命周期契约。
func (*Exporter) Shutdown(context.Context) error { return nil }

// ExportSpans 输出 span 快照。
func (e *Exporter) ExportSpans(_ context.Context, spans []observe.SpanSnapshot) error {
	for _, span := range spans {
		if err := e.write(record{
			Signal:       "trace",
			Name:         span.Name,
			Resource:     span.Resource.Name,
			Scope:        span.Scope.Name,
			TraceID:      span.SpanContext.TraceID.String(),
			SpanID:       span.SpanContext.SpanID.String(),
			ParentSpanID: validSpanID(span.Parent.SpanID),
			Attributes:   attributes(span.Attributes),
		}); err != nil {
			return err
		}
	}
	return nil
}

// ExportMetrics 输出指标快照摘要。
func (e *Exporter) ExportMetrics(_ context.Context, metrics []observe.MetricData) error {
	for _, metric := range metrics {
		if err := e.write(record{
			Signal:   "metric",
			Name:     metric.Name,
			Resource: metric.Resource.Name,
			Scope:    metric.Scope.Name,
			Points:   pointCount(metric.Aggregation),
		}); err != nil {
			return err
		}
	}
	return nil
}

// ExportLogs 输出结构化日志摘要。
func (e *Exporter) ExportLogs(_ context.Context, logs []observe.LogRecord) error {
	for _, item := range logs {
		if err := e.write(record{
			Signal:   "log",
			Name:     "log.record",
			Resource: item.Resource.Name,
			Scope:    item.Scope.Name,
			TraceID:  validTraceID(item.SpanContext.TraceID),
			SpanID:   validSpanID(item.SpanContext.SpanID),
			Level:    item.Level.String(),
			Message:  item.Message,
		}); err != nil {
			return err
		}
	}
	return nil
}

// ExportEvents 输出离散事件摘要。
func (e *Exporter) ExportEvents(_ context.Context, events []observe.EventRecord) error {
	for _, event := range events {
		if err := e.write(record{
			Signal:   "event",
			Name:     event.Name,
			Resource: event.Resource.Name,
			Scope:    event.Scope.Name,
			TraceID:  validTraceID(event.SpanContext.TraceID),
			SpanID:   validSpanID(event.SpanContext.SpanID),
			Value:    event.Body.String(),
		}); err != nil {
			return err
		}
	}
	return nil
}

func (e *Exporter) write(value record) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.encoder.Encode(value)
}

func pointCount(aggregation observe.MetricAggregation) int {
	switch data := aggregation.(type) {
	case observe.GaugeData:
		return len(data.Points)
	case *observe.GaugeData:
		return len(data.Points)
	case observe.SumData:
		return len(data.Points)
	case *observe.SumData:
		return len(data.Points)
	case observe.HistogramData:
		return len(data.Points)
	case *observe.HistogramData:
		return len(data.Points)
	default:
		return 0
	}
}

func validTraceID(id observe.TraceID) string {
	if !id.IsValid() {
		return ""
	}
	return id.String()
}

func validSpanID(id observe.SpanID) string {
	if !id.IsValid() {
		return ""
	}
	return id.String()
}

func attributes(attrs []observe.Attr) map[string]string {
	if len(attrs) == 0 {
		return nil
	}
	values := make(map[string]string, len(attrs))
	for _, attr := range attrs {
		values[attr.Key.String()] = attr.Value.String()
	}
	return values
}
