package zaplog

import (
	"context"
	"testing"
)

func TestDevelopmentLogger(t *testing.T) {
	logger := StdLogger().WithArgs("function", "TestDevelopmentLogger")

	ctx := context.Background()
	logger.WithContext(ctx).Debug("hello")

	StdContextHook = func(ctx context.Context) []interface{} {
		return []interface{}{"TraceID", "123456"}
	}
	logger.WithContext(ctx).Debug("hello")
	logger.WithContext(ctx).Warn("warn")
	logger.WithContext(ctx).Error("error")
}
