package logger

import (
	"context"

	"go.uber.org/zap"
)

var (
	// Key that can be used to get the logger from the request context.
	loggerKey = "logger"
)

func FromContext(ctx context.Context) (*zap.Logger, bool) {
	l, ok := ctx.Value(loggerKey).(*zap.Logger)
	return l, ok
}

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
