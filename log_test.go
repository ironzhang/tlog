package tlog_test

import (
	"context"
	"fmt"
	"os"

	"github.com/ironzhang/tlog"
	"github.com/ironzhang/tlog/iface"
	"github.com/ironzhang/tlog/zaplog"
)

func ContextTestHook(ctx context.Context) []interface{} {
	return []interface{}{"TraceID", "123456"}
}

func RecoverPanic(f func()) {
	defer func() {
		recover()
	}()
	f()
}

func PrintLog(name string) {
	min, max := iface.DEBUG, iface.PANIC
	for lv := min; lv <= max; lv++ {
		RecoverPanic(func() {
			tlog.Print(0, lv, "Print ", lv, " ", name)
		})
		RecoverPanic(func() {
			tlog.Printf(0, lv, "Printf %s %s", lv, name)
		})
		RecoverPanic(func() {
			tlog.Printw(0, lv, "Printw", "level", lv, "name", name)
		})
	}

	tlog.Debug("Debug ", name)
	tlog.Debugf("Debugf %s", name)
	tlog.Debugw("Debugw", "name", name)

	tlog.Info("Info ", name)
	tlog.Infof("Infof %s", name)
	tlog.Infow("Infow", "name", name)

	tlog.Warn("Warn ", name)
	tlog.Warnf("Warnf %s", name)
	tlog.Warnw("Warnw", "name", name)

	tlog.Error("Error ", name)
	tlog.Errorf("Errorf %s", name)
	tlog.Errorw("Errorw", "name", name)

	RecoverPanic(func() {
		tlog.Panic("Panic", name)
	})
	RecoverPanic(func() {
		tlog.Panicf("Panicf %s", name)
	})
	RecoverPanic(func() {
		tlog.Panicw("Panicw", "name", name)
	})

	tlog.WithArgs("name", name).Info("with args")
	tlog.WithContext(context.Background()).Info("with context")
	tlog.WithArgs("name", name).WithContext(context.Background()).Info("with args and with context")
}

func ExampleNopLogger() {
	tlog.SetFactory(nil)
	PrintLog("ExampleNopLogger")

	// output:
}

func ExampleStdLogger() {
	zaplog.StdContextHook = ContextTestHook
	PrintLog("ExampleStdLogger")

	// output:
}

func ExampleLogger() {
	logger, err := zaplog.New(zaplog.NewDevelopmentConfig(), zaplog.SetContextHook(zaplog.ContextHookFunc(ContextTestHook)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "new logger: %v", err)
		return
	}
	tlog.SetFactory(logger)
	PrintLog("ExampleLogger")

	// output:
}
