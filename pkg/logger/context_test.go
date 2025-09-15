package logger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithLogger(t *testing.T) {
	expectedLogger := Setup()

	ctx := WithLogger(context.Background(), expectedLogger)
	got := ctx.Value(loggerKey)

	assert.Equal(t, expectedLogger, got)
}

func TestFromContext_OK(t *testing.T) {
    expectedLogger := Setup()

    ctx := WithLogger(context.Background(), expectedLogger)
    got := FromContext(ctx)

    assert.Equal(t, expectedLogger, got)
}
