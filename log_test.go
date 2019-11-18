package tlog_test

import (
	"context"
	"testing"

	"github.com/ironzhang/tlog"
	"github.com/ironzhang/tlog/zaplog"
)

func TestLog(t *testing.T) {
	zaplog.StdContextHook = func(ctx context.Context) []interface{} {
		return []interface{}{"TraceID", "123456"}
	}

	tlog.Debug("hello")
	tlog.Info("hello")
	tlog.WithArgs("function", "TestLog").Warn("hello")
	tlog.WithContext(context.Background()).WithArgs("function", "TestLog").Error("hello")
}
