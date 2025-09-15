package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestSetup exists only for 100% test coverage.
func TestSetup(t *testing.T) {
	got := Setup()

	assert.Equal(t, zap.InfoLevel, got.Level())
	assert.Equal(t, "json", zap.NewProductionConfig().Encoding)
}
