package logger

import (
	"context"

	"github.com/opentracing/opentracing-go"
	opentracinglog "github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Reporter ...
var logger *zap.Logger

// Init call
func Init(filepath string) {
	cfg := zap.NewProductionConfig()

	if filepath != "" {
		cfg.OutputPaths = []string{filepath}
	} else {
		cfg.OutputPaths = []string{"stderr"}
	}

	logger, _ = cfg.Build()
}

// InfoWithContext ...
func InfoWithContext(ctx context.Context, log string, fields ...zapcore.Field) {
	InfoWithSpan(opentracing.SpanFromContext(ctx), log, fields...)
}

// InfoWithSpan ...
func InfoWithSpan(span opentracing.Span, log string, fields ...zapcore.Field) {
	traceID := span.Context().(jaeger.SpanContext).TraceID()
	fields = append(fields, zap.String("trace_id", traceID.String()))

	Info(log, fields...)
	logSpan(span, log, fields...)
}

// Info ...
func Info(log string, fields ...zapcore.Field) {
	logger.Info(log, fields...)
}

// ErrorWithContext ...
func ErrorWithContext(ctx context.Context, log string, fields ...zapcore.Field) {
	ErrorWithSpan(opentracing.SpanFromContext(ctx), log, fields...)
}

// ErrorWithSpan ...
func ErrorWithSpan(span opentracing.Span, log string, fields ...zapcore.Field) {
	traceID := span.Context().(jaeger.SpanContext).TraceID()
	fields = append(fields, zap.String("trace_id", traceID.String()))

	Error(log, fields...)
	logSpan(span, log, fields...)
}

// Error ...
func Error(log string, fields ...zapcore.Field) {
	logger.Error(log, fields...)
}

func logSpan(span opentracing.Span, log string, fields ...zapcore.Field) {
	if span == nil {
		return
	}

	of := []opentracinglog.Field{}
	if log != "" {
		of = append(of, opentracinglog.String("event", log))
	}

	if len(fields) > 0 {
		of = append(of, ZapFieldsToOpentracing(fields...)...)
	}

	span.LogFields(of...)
}
