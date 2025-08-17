package logger

import "go.uber.org/zap"

// Setup creates new logger instance.
// TODO: add more options
func Setup() *zap.Logger {
	logger := zap.Must(zap.NewProduction())
	return logger
}
