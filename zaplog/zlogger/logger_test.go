package zlogger

import (
	"context"
	"os"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewTestCore() zapcore.Core {
	enc := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	return zapcore.NewCore(enc, os.Stdout, zapcore.DebugLevel)
}

func TestLoggerPrint(t *testing.T) {
	core := NewTestCore()
	log := New("", DEBUG, core, zap.AddCaller())

	min, max := DEBUG, ERROR
	for lvl := min; lvl <= max; lvl++ {
		log.Print(0, lvl, "Print level", lvl)
		log.Printf(0, lvl, "Printf level=%d", lvl)
		log.Printw(0, lvl, "Printw", "level", lvl)
	}
}

func TestLogger(t *testing.T) {
	core := NewTestCore()
	logger := New("", DEBUG, core, zap.AddCaller())
	logger.SetWithContextFunc(func(ctx context.Context) []interface{} {
		return []interface{}{"trace_id", "123456"}
	})
	l := logger.WithArgs("function", "TestLogger").WithContext(context.Background())

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
