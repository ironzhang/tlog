package logger

import (
	"context"
)

// Logger 日志接口
type Logger interface {
	Named(name string) Logger
	WithArgs(args ...interface{}) Logger
	WithContext(ctx context.Context) Logger

	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
	Tracew(message string, kvs ...interface{})

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
}
