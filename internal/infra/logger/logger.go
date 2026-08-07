package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"go.opentelemetry.io/otel/trace"
)

// LevelCritical Custom levels
const LevelCritical = slog.LevelError + 1

const (
	timeFormat = "01-02-2006 15:04:05.000"
)

// Logger encapsulates handlers for different environments.
type Logger struct {
	json *slog.Logger
	text *slog.Logger
}

// NewLogger builds the handlers and sets the global default.
func NewLogger() *Logger {
	// Options for JSON (Production)
	jsonOpts := &slog.HandlerOptions{
		Level:       slog.LevelDebug,
		ReplaceAttr: replaceAttrJSON,
	}

	// Options for Text (Development)
	textOpts := &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: timeFormat,
		NoColor:    false,
	}

	l := &Logger{
		json: slog.New(slog.NewJSONHandler(os.Stdout, jsonOpts)),
		text: slog.New(tint.NewHandler(os.Stderr, textOpts)),
	}

	// Sets the system default based on the environment
	if os.Getenv("ENV") == "development" {
		slog.SetDefault(l.text)
	} else {
		slog.SetDefault(l.json)
	}

	return l
}

// ── JSON Methods (Production / Loki) ───────────────────────────────────────────

func (l *Logger) DebugJSON(msg string, args ...any) { l.json.Debug(msg, args...) }
func (l *Logger) InfoJSON(msg string, args ...any)  { l.json.Info(msg, args...) }
func (l *Logger) WarnJSON(msg string, args ...any)  { l.json.Warn(msg, args...) }
func (l *Logger) ErrorJSON(msg string, args ...any) { l.json.Error(msg, args...) }
func (l *Logger) CriticalJSON(msg string, args ...any) {
	l.json.Log(context.Background(), LevelCritical, msg, args...)
}

// ── Text Methods (Development / Terminal) ────────────────────────────────

func (l *Logger) DebugText(msg string, args ...any) { l.text.Debug(msg, args...) }
func (l *Logger) InfoText(msg string, args ...any)  { l.text.Info(msg, args...) }
func (l *Logger) WarnText(msg string, args ...any)  { l.text.Warn(msg, args...) }
func (l *Logger) ErrorText(msg string, args ...any) { l.text.Error(msg, args...) }
func (l *Logger) CriticalText(msg string, args ...any) {
	l.text.Log(context.Background(), LevelCritical, msg, args...)
}

// ── Observability ───────────────────────────────────────────────────────────

// WithTrace returns the JSON logger enriched with OpenTelemetry data.
func (l *Logger) WithTrace(ctx context.Context) *slog.Logger {
	span := trace.SpanContextFromContext(ctx)
	if !span.IsValid() || !span.HasTraceID() {
		return l.json
	}
	return l.json.With(
		slog.String("trace_id", span.TraceID().String()),
		slog.String("span_id", span.SpanID().String()),
	)
}

// ── Helpers and Attributes ───────────────────────────────────────────────────────

func replaceAttrJSON(_ []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.TimeKey:
		// Optimization: uses the time that slog already provides
		a.Value = slog.StringValue(a.Value.Time().Format(timeFormat))
	case slog.LevelKey:
		a.Value = slog.StringValue(levelLabel(a))
	}
	return a
}

func levelLabel(a slog.Attr) string {
	level, ok := a.Value.Any().(slog.Level)
	if !ok {
		return a.Value.String()
	}
	switch level {
	case LevelCritical:
		return "CRITICAL"
	case slog.LevelDebug:
		return "DEBUG"
	case slog.LevelInfo:
		return "INFO"
	case slog.LevelWarn:
		return "WARN"
	case slog.LevelError:
		return "ERROR"
	default:
		return level.String()
	}
}

// SlogJSON and SlogText allow direct access if necessary
func (l *Logger) SlogJSON() *slog.Logger { return l.json }
func (l *Logger) SlogText() *slog.Logger { return l.text }
