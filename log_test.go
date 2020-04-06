package tlog_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"git.xiaojukeji.com/pearls/tlog"
	"git.xiaojukeji.com/pearls/tlog/iface"
	"git.xiaojukeji.com/pearls/tlog/zaplog"
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

func PrintAccessLogs(name string) {
	logger := tlog.Named("access").WithArgs("logger", "access")

	logger.Debug("Debug ", name)
	logger.Debugf("Debugf %s", name)
	logger.Debugw("Debugw", "name", name)

	logger.Info("Info ", name)
	logger.Infof("Infof %s", name)
	logger.Infow("Infow", "name", name)

	logger.Warn("Warn ", name)
	logger.Warnf("Warnf %s", name)
	logger.Warnw("Warnw", "name", name)

	logger.Error("Error ", name)
	logger.Errorf("Errorf %s", name)
	logger.Errorw("Errorw", "name", name)

	RecoverPanic(func() {
		logger.Panic("Panic ", name)
	})
	RecoverPanic(func() {
		logger.Panicf("Panicf %s", name)
	})
	RecoverPanic(func() {
		logger.Panicw("Panicw", "name", name)
	})

	logger.Print(0, iface.INFO, "Print ", name)
	logger.Printf(0, iface.INFO, "Printf %s", name)
	logger.Printw(0, iface.INFO, "Printw", "name", name)

	logger.WithArgs("name", name).Info("with args")
	logger.WithContext(context.Background()).Info("with context")
	logger.WithArgs("name", name).WithContext(context.Background()).Info("with args and with context")
}

func PrintLogs(name string) {
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
		tlog.Panic("Panic ", name)
	})
	RecoverPanic(func() {
		tlog.Panicf("Panicf %s", name)
	})
	RecoverPanic(func() {
		tlog.Panicw("Panicw", "name", name)
	})

	tlog.Print(0, iface.INFO, "Print ", name)
	tlog.Printf(0, iface.INFO, "Printf %s", name)
	tlog.Printw(0, iface.INFO, "Printw", "name", name)

	tlog.WithArgs("name", name).Info("with args")
	tlog.WithContext(context.Background()).Info("with context")
	tlog.WithArgs("name", name).WithContext(context.Background()).Info("with args and with context")

	PrintAccessLogs(name)
}

func TestStdLogger(t *testing.T) {
	zaplog.StdContextHook = ContextTestHook
	PrintLogs("TestStdLogger")
}

func TestNopLogger(t *testing.T) {
	tlog.SetLogger(nil)
	PrintLogs("TestNopLogger")
	tlog.SetLogger(zaplog.StdLogger())
}

func TestLogger(t *testing.T) {
	logger, err := zaplog.New(zaplog.NewDevelopmentConfig(), zaplog.SetContextHook(zaplog.ContextHookFunc(ContextTestHook)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "new logger: %v", err)
		return
	}
	tlog.SetLogger(logger)
	PrintLogs("ExampleLogger")
	tlog.SetLogger(zaplog.StdLogger())
}
