package utils

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

func TraceIdFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}

	return ""
}

func CopyContext(ctx context.Context) context.Context {
	span := trace.SpanContextFromContext(ctx)
	return trace.ContextWithSpanContext(context.Background(), span)
}
