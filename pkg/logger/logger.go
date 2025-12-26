package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

type Logger struct {
	*slog.Logger
}

type Config struct {
	Level      string // debug, info, warn, error
	Format     string // json, text
	Output     io.Writer
	AddSource  bool // file/line
	TimeFormat string
}

type ContextKey string

const (
	RequestIDKey ContextKey = "request_id"
	UserIDKey    ContextKey = "user_id"
	TraceIDKey   ContextKey = "trace_id"
)

var defaultLogger *Logger

func New(cfg Config) *Logger {
	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}

	if cfg.TimeFormat == "" {
		cfg.TimeFormat = time.RFC3339
	}

	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: cfg.AddSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format(cfg.TimeFormat))
				}
			}
			if a.Key == slog.SourceKey {
				if source, ok := a.Value.Any().(*slog.Source); ok {
					if idx := strings.Index(source.File, "nexus/"); idx != -1 {
						source.File = source.File[idx:]
					}
				}
			}
			return a
		},
	}

	var handler slog.Handler
	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(cfg.Output, opts)
	} else {
		handler = slog.NewTextHandler(cfg.Output, opts)
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}

func Init(cfg Config) {
	defaultLogger = New(cfg)
	slog.SetDefault(defaultLogger.Logger)
}

func Default() *Logger {
	if defaultLogger == nil {
		defaultLogger = New(Config{
			Level:  "info",
			Format: "text",
			Output: os.Stdout,
		})
	}
	return defaultLogger
}

func FromContext(ctx context.Context) *Logger {
	logger := Default()

	attrs := make([]any, 0, 6)

	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		attrs = append(attrs, slog.String("request_id", requestID.(string)))
	}

	if userID := ctx.Value(UserIDKey); userID != nil {
		attrs = append(attrs, slog.String("user_id", userID.(string)))
	}

	if traceID := ctx.Value(TraceIDKey); traceID != nil {
		attrs = append(attrs, slog.String("trace_id", traceID.(string)))
	}

	if len(attrs) > 0 {
		return &Logger{
			Logger: logger.With(attrs...),
		}
	}

	return logger
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	attrs := make([]any, 0, 6)

	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		attrs = append(attrs, slog.String("request_id", requestID.(string)))
	}

	if userID := ctx.Value(UserIDKey); userID != nil {
		attrs = append(attrs, slog.String("user_id", userID.(string)))
	}

	if traceID := ctx.Value(TraceIDKey); traceID != nil {
		attrs = append(attrs, slog.String("trace_id", traceID.(string)))
	}

	if len(attrs) > 0 {
		return &Logger{
			Logger: l.With(attrs...),
		}
	}

	return l
}

func (l *Logger) WithFields(fields map[string]any) *Logger {
	attrs := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		attrs = append(attrs, slog.Any(k, v))
	}
	return &Logger{
		Logger: l.With(attrs...),
	}
}

func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		Logger: l.With(slog.String("error", err.Error())),
	}
}

func Debug(msg string, args ...any) {
	Default().Debug(msg, args...)
}

func Info(msg string, args ...any) {
	Default().Info(msg, args...)
}

func Warn(msg string, args ...any) {
	Default().Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Default().Error(msg, args...)
}

func DebugContext(ctx context.Context, msg string, args ...any) {
	FromContext(ctx).Debug(msg, args...)
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	FromContext(ctx).Info(msg, args...)
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	FromContext(ctx).Warn(msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	FromContext(ctx).Error(msg, args...)
}

func Fatal(msg string, args ...any) {
	Default().Error(msg, args...)
	os.Exit(1)
}

func FatalContext(ctx context.Context, msg string, args ...any) {
	FromContext(ctx).Error(msg, args...)
	os.Exit(1)
}
