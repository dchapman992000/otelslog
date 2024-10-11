package otelslog

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	"go.opentelemetry.io/otel/trace"
)

type ContextHandler struct {
	slog.Handler
}

var DataDogFields bool = false

// Taken from https://darrenparkinson.uk/posts/2023-09-14-datadog-log-correlation-with-slog/
func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(h.addTraceFromContext(ctx)...)
	return h.Handler.Handle(ctx, r)
}

func (h ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return ContextHandler{h.Handler.WithAttrs(attrs)}
}

func (h ContextHandler) WithGroup(name string) slog.Handler {
	return ContextHandler{h.Handler.WithGroup(name)}
}

func (h ContextHandler) addTraceFromContext(ctx context.Context) (as []slog.Attr) {
	span := trace.SpanFromContext(ctx)
	if span != nil && DataDogFields {
		traceID := convertTraceID(span.SpanContext().TraceID().String())
		spanID := convertTraceID(span.SpanContext().SpanID().String())

		ddgroup := slog.Group("dd",
			slog.String("trace_id", traceID),
			slog.String("span_id", spanID),
			slog.String("env", os.Getenv("OTEL_RESOURCE_ATTRIBUTES")),
			slog.String("service", os.Getenv("OTEL_SERVICE_NAME")),
			slog.String("version", "1.0"),
			slog.String("source", "go"))

		as = append(as, ddgroup)
	}
	return
}

func InitialiseLogging(handy slog.Handler) *slog.Logger {
	ctxHandler := ContextHandler{handy}

	// WithGroup duplicates or overwrites these - there has to be a better way
	/* ddvars := slog.Group("dd",
		slog.String("env", os.Getenv("OTEL_RESOURCE_ATTRIBUTES")),
		slog.String("service", os.Getenv("OTEL_SERVICE_NAME")),
		slog.String("version", "1.0"),
		slog.String("source", "go"),

	)
	logger := slog.New(ctxHandler).With(ddvars) */

	logger := slog.New(ctxHandler)
	slog.SetDefault(logger)
	return logger
}

// Taken from https://docs.datadoghq.com/tracing/other_telemetry/connect_logs_and_traces/opentelemetry/?tab=go
func convertTraceID(id string) string {
	if len(id) < 16 {
		return ""
	}
	if len(id) > 16 {
		id = id[16:]
	}
	intValue, err := strconv.ParseUint(id, 16, 64)
	if err != nil {
		return ""
	}
	return strconv.FormatUint(intValue, 10)
}
