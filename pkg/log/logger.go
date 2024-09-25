// Package log provides context-aware and structured logging capabilities.
package log

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Logger is a logger that supports log levels, context and structured logging.
type Logger interface {
	// With returns a logger based off the root logger and decorates it with the given context and arguments.
	With(ctx context.Context, args ...interface{}) Logger

	// Debug uses fmt.Sprint to construct and log a message at DEBUG level
	Debug(args ...interface{})
	// Info uses fmt.Sprint to construct and log a message at INFO level
	Info(args ...interface{})
	// Error uses fmt.Sprint to construct and log a message at ERROR level
	Error(args ...interface{})

	// Debugf uses fmt.Sprintf to construct and log a message at DEBUG level
	Debugf(format string, args ...interface{})
	// Infof uses fmt.Sprintf to construct and log a message at INFO level
	Infof(format string, args ...interface{})
	// Errorf uses fmt.Sprintf to construct and log a message at ERROR level
	Errorf(format string, args ...interface{})
}

type logger struct {
	*zap.SugaredLogger
}

type contextKey int

const (
	requestIDKey contextKey = iota
	correlationIDKey
)

// New creates a new logger using the default configuration.
func New() Logger {
	l, _ := zap.NewProduction()
	return NewWithZap(l)
}

// NewWithZap creates a new logger using the preconfigured zap logger.
func NewWithZap(l *zap.Logger) Logger {
	return &logger{l.Sugar()}
}

func (l *logger) With(ctx context.Context, args ...interface{}) Logger {
	if ctx != nil {
		if id, ok := ctx.Value(requestIDKey).(string); ok {
			args = append(args, zap.String("request_id", id))
		}
		if id, ok := ctx.Value(correlationIDKey).(string); ok {
			args = append(args, zap.String("correlation_id", id))
		}
	}
	if len(args) > 0 {
		return &logger{l.SugaredLogger.With(args...)}
	}
	return l
}

func getRequestId(req *http.Request) string {
	return req.Header.Get("X-Request-ID")
}

func getCorrelationId(req *http.Request) string {
	return req.Header.Get("X-Correlation-ID")
}

func WithRequest(ctx context.Context, req *http.Request) context.Context {
	id := getRequestId(req)

	if id == "" {
		id = uuid.New().String()
	}

	ctx = context.WithValue(ctx, requestIDKey, id)
	correlationId := getCorrelationId(req)
	if correlationId != "" {
		ctx = context.WithValue(ctx, correlationIDKey, correlationId)
	}

	return ctx
}
