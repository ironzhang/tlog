package tlog_test

import (
	"context"
	"fmt"
	"os"

	"github.com/ironzhang/tlog"
	"github.com/ironzhang/tlog/zaplog"
)

func ExampleStdLogger() {
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

func ExampleLogger() {
	logger, err := zaplog.New(zaplog.NewDevelopmentConfig())
	if err != nil {
		fmt.Fprintf(os.Stderr, "new logger: %v", err)
		return
	}
	tlog.SetFactory(logger)

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
