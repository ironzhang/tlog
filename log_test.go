package tlog_test

import (
	"context"
	"testing"

	"github.com/ironzhang/tlog"
	"github.com/ironzhang/tlog/zaplog"
)

func TestLog(t *testing.T) {
	//	w := &lumberjack.Logger{
	//		Filename: "./test.log",
	//	}
	//	defer w.Close()

	hook := func(ctx context.Context) []interface{} {
		return []interface{}{"TraceID", "123456"}
	}
	l := zaplog.Std.WithOptions(zaplog.SetContextHook(hook))
	tlog.SetLogger(l)

	tlog.Trace("hello")
	tlog.Debug("hello")
	tlog.Info("hello")
	tlog.WithArgs("function", "TestLog").Warn("hello")
	tlog.WithContext(context.Background()).WithArgs("function", "TestLog").Warn("hello")
}
