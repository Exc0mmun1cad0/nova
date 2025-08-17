package logger

import (
	"context"

	"go.uber.org/zap"
)

var (
	// Key that can be used to get the logger from the request context.
	loggerKey = "logger"
)

func FromContext(ctx context.Context) *zap.Logger {
	l, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		panic("object of wrong type is available via logger key")
	}
	return l
}

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
