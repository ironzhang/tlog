package zaplog

import (
	"context"
	"testing"
)

func TestDevelopmentLogger(t *testing.T) {
	logger := DevelopmentLogger.WithArgs("function", "TestDevelopmentLogger")

	ctx := context.Background()
	logger.WithContext(ctx).Debug("hello")

	DevelopmentContextHook = func(ctx context.Context) []interface{} {
		return []interface{}{"TraceID", "123456"}
	}
	logger.WithContext(ctx).Debug("hello")
}
