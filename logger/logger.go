package logger

import "context"

type Logger interface {
	WithArgs(args ...interface{}) Logger
	WithContext(ctx context.Context) Logger

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugw(message string, kvs ...interface{})

	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
	Tracew(message string, kvs ...interface{})

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

var WithContextHook = func(ctx context.Context) []interface{} {
	return nil
}
