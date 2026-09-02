// Package example 展示 observe-sdk 四类信号及其生命周期。
package example

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/goark-projects/goark-observe-sdk-example/internal/consoleexporter"
	"goark.dev/observe"
	observesdk "goark.dev/observe-sdk"
)

const scopeName = "example/orders"

// Run 生成具有关联上下文的 trace、metric、log 和 event。
func Run(ctx context.Context, output io.Writer) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	provider, err := observesdk.NewProvider(
		observesdk.WithResource(observe.NewResource(
			"observe-sdk-example",
			observe.WithResourceVersion("1.0.0"),
			observe.WithResourceEnv("development"),
		)),
		observesdk.WithExporters(consoleexporter.New(output)),
		observesdk.WithMetricCardinalityLimit(128),
	)
	if err != nil {
		return err
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = errors.Join(err, provider.Shutdown(shutdownCtx))
	}()

	meter := provider.Meter(scopeName)
	requests, err := meter.Int64Counter(
		"orders.processed",
		observe.WithDescription("Processed orders"),
		observe.WithUnit(observe.UnitCount),
		observe.WithAttributeKeys("outcome"),
	)
	if err != nil {
		return err
	}
	duration, err := meter.Float64Histogram(
		"orders.processing.duration",
		observe.WithDescription("Order processing duration"),
		observe.WithUnit(observe.UnitSeconds),
		observe.WithAttributeKeys("outcome"),
		observe.WithExplicitBounds(.005, .01, .025, .05, .1, .25, .5),
	)
	if err != nil {
		return err
	}

	started := time.Now()
	ctx, span := provider.Tracer(scopeName).Start(
		ctx,
		"orders.process",
		observe.WithSpanKind(observe.SpanKindInternal),
		observe.WithSpanAttrs(observe.String("order.type", "standard")),
	)
	provider.Logger(scopeName).Log(
		ctx,
		observe.LevelInfo,
		"order processing started",
		observe.String("order.type", "standard"),
	)
	provider.Eventer(scopeName).Emit(
		ctx,
		"order.accepted",
		observe.String("channel", "example"),
	)
	span.AddEvent("inventory.reserved", observe.Int("items", 2))
	requests.Add(ctx, 1, observe.String("outcome", observe.OutcomeOK))
	duration.Record(ctx, time.Since(started).Seconds(), observe.String("outcome", observe.OutcomeOK))
	span.SetStatus(observe.StatusOK, "")
	span.End()

	return provider.ForceFlush(ctx)
}
