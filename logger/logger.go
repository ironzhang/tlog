package logger

import (
	"context"
)

// Logger 日志接口
type Logger interface {
	WithArgs(args ...interface{}) Logger
	WithContext(ctx context.Context) Logger

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugw(message string, kvs ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infow(message string, kvs ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Warnw(message string, kvs ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorw(message string, kvs ...interface{})

	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Panicw(message string, kvs ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalw(message string, kvs ...interface{})

	Print(depth int, level Level, args ...interface{})
	Printf(depth int, level Level, format string, args ...interface{})
	Printw(depth int, level Level, message string, kvs ...interface{})
}
