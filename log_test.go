package tlog_test

import (
	"context"

	"github.com/ironzhang/tlog"
	"github.com/ironzhang/tlog/zaplog"
)

func ExampleLog() {
	zaplog.StdContextHook = func(ctx context.Context) []interface{} {
		return []interface{}{"TraceID", "123456"}
	}

	tlog.Debug("debug")
	tlog.Info("info")
	tlog.Warn("warn")
	tlog.Error("error")
	func() {
		defer func() {
			recover()
		}()
		tlog.Panic("panic")
	}()

	tlog.WithArgs("function", "TestLog").Info("hello")
	tlog.WithContext(context.Background()).WithArgs("function", "TestLog").Info("hello")

	// output:
}
