package zlogger

import (
	"context"
	"os"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ironzhang/tlog/iface"
)

type TestContextHook struct{}

func (p TestContextHook) WithContext(ctx context.Context) []interface{} {
	return []interface{}{"trace_id", "123456"}
}

func NewTestCore() zapcore.Core {
	enc := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	return zapcore.NewCore(enc, os.Stdout, zapcore.DebugLevel)
}

func NewTestLogger() *Logger {
	core := NewTestCore()
	return New("", core, TestContextHook{}, zap.AddCaller())
}

func TestLoggerPrint(t *testing.T) {
	logger := NewTestLogger()
	min, max := iface.DEBUG, iface.ERROR
	for lvl := min; lvl <= max; lvl++ {
		logger.Print(0, lvl, "Print level", lvl)
		logger.Printf(0, lvl, "Printf level=%s", lvl)
		logger.Printw(0, lvl, "Printw", "level", lvl)
	}
}

func TestLogger(t *testing.T) {
	logger := NewTestLogger().WithArgs("function", "TestLogger").WithContext(context.Background())

	type LogFunc func(args ...interface{})
	logFuncs := []LogFunc{
		logger.Debug,
		logger.Info,
		logger.Warn,
		logger.Error,
	}
	for _, log := range logFuncs {
		log("hello", "world")
	}

	type LogfFunc func(format string, args ...interface{})
	logfFuncs := []LogfFunc{
		logger.Debugf,
		logger.Infof,
		logger.Warnf,
		logger.Errorf,
	}
	for _, log := range logfFuncs {
		log("hello, %s", "world")
	}

	type LogwFunc func(message string, kvs ...interface{})
	logwFuncs := []LogwFunc{
		logger.Debugw,
		logger.Infow,
		logger.Warnw,
		logger.Errorw,
	}
	for _, log := range logwFuncs {
		log("hello, world", "hello", "world")
	}
}
