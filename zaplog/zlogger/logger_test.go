package zlogger

import (
	"context"
	"os"
	"testing"

	"github.com/ironzhang/tlog/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	log := NewTestLogger()
	min, max := logger.DEBUG, logger.ERROR
	for lvl := min; lvl <= max; lvl++ {
		log.Print(0, lvl, "Print level", lvl)
		log.Printf(0, lvl, "Printf level=%d", lvl)
		log.Printw(0, lvl, "Printw", "level", lvl)
	}
}

func TestLogger(t *testing.T) {
	log := NewTestLogger()
	l := log.WithArgs("function", "TestLogger").WithContext(context.Background())

	type LogFunc func(args ...interface{})
	logFuncs := []LogFunc{
		l.Debug,
		l.Info,
		l.Warn,
	}
	for _, log := range logFuncs {
		log("hello", "world")
	}

	type LogfFunc func(format string, args ...interface{})
	logfFuncs := []LogfFunc{
		l.Debugf,
		l.Infof,
	}
	for _, log := range logfFuncs {
		log("hello, %s", "world")
	}

	type LogwFunc func(message string, kvs ...interface{})
	logwFuncs := []LogwFunc{
		l.Debugw,
		l.Infow,
	}
	for _, log := range logwFuncs {
		log("hello, world", "hello", "world")
	}
}
