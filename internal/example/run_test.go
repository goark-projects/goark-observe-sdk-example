package example_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bytedance/sonic"

	"github.com/goark-projects/goark-observe-sdk-example/internal/example"
)

func TestRunExportsAllSignalsWithTraceCorrelation(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	if err := example.Run(t.Context(), &output); err != nil {
		t.Fatalf("Run: %v", err)
	}

	type record struct {
		Signal  string `json:"signal"`
		TraceID string `json:"trace_id"`
	}
	counts := make(map[string]int)
	var traceID string
	var correlatedTraceIDs []string
	for _, line := range strings.Split(strings.TrimSpace(output.String()), "\n") {
		var item record
		if err := sonic.Unmarshal([]byte(line), &item); err != nil {
			t.Fatalf("invalid JSON line %q: %v", line, err)
		}
		counts[item.Signal]++
		if item.Signal == "trace" {
			traceID = item.TraceID
		}
		if item.Signal == "log" || item.Signal == "event" {
			correlatedTraceIDs = append(correlatedTraceIDs, item.TraceID)
		}
	}
	for _, signal := range []string{"trace", "metric", "log", "event"} {
		if counts[signal] == 0 {
			t.Fatalf("signal %q was not exported: %v", signal, counts)
		}
	}
	if traceID == "" {
		t.Fatal("trace ID was not exported")
	}
	for _, correlated := range correlatedTraceIDs {
		if correlated != traceID {
			t.Fatalf("correlated trace ID = %q, want %q", correlated, traceID)
		}
	}
}
