package tlog

import (
	"context"

	"github.com/ironzhang/tlog/iface"
	"github.com/ironzhang/tlog/zaplog"
)

type Level = iface.Level

const (
	DEBUG = iface.DEBUG
	INFO  = iface.INFO
	WARN  = iface.WARN
	ERROR = iface.ERROR
	PANIC = iface.PANIC
	FATAL = iface.FATAL
)

type Logger = iface.Logger

type nopLogger struct {
}

func (p nopLogger) Named(name string) Logger                                          { return p }
func (p nopLogger) WithArgs(args ...interface{}) Logger                               { return p }
func (p nopLogger) WithContext(ctx context.Context) Logger                            { return p }
func (p nopLogger) Debug(args ...interface{})                                         {}
func (p nopLogger) Debugf(format string, args ...interface{})                         {}
func (p nopLogger) Debugw(message string, kvs ...interface{})                         {}
func (p nopLogger) Info(args ...interface{})                                          {}
func (p nopLogger) Infof(format string, args ...interface{})                          {}
func (p nopLogger) Infow(message string, kvs ...interface{})                          {}
func (p nopLogger) Warn(args ...interface{})                                          {}
func (p nopLogger) Warnf(format string, args ...interface{})                          {}
func (p nopLogger) Warnw(message string, kvs ...interface{})                          {}
func (p nopLogger) Error(args ...interface{})                                         {}
func (p nopLogger) Errorf(format string, args ...interface{})                         {}
func (p nopLogger) Errorw(message string, kvs ...interface{})                         {}
func (p nopLogger) Panic(args ...interface{})                                         {}
func (p nopLogger) Panicf(format string, args ...interface{})                         {}
func (p nopLogger) Panicw(message string, kvs ...interface{})                         {}
func (p nopLogger) Fatal(args ...interface{})                                         {}
func (p nopLogger) Fatalf(format string, args ...interface{})                         {}
func (p nopLogger) Fatalw(message string, kvs ...interface{})                         {}
func (p nopLogger) Print(depth int, level Level, args ...interface{})                 {}
func (p nopLogger) Printf(depth int, level Level, format string, args ...interface{}) {}
func (p nopLogger) Printw(depth int, level Level, message string, kvs ...interface{}) {}

func init() {
	SetLogger(zaplog.StdLogger())
}

var logging Logger

func SetLogger(l Logger) (old Logger) {
	old = logging
	if l == nil {
		l = nopLogger{}
	}
	logging = l
	return old
}

func GetLogger() Logger {
	return logging
}

func Named(name string) Logger {
	return logging.Named(name)
}

func WithArgs(args ...interface{}) Logger {
	return logging.WithArgs(args...)
}

func WithContext(ctx context.Context) Logger {
	return logging.WithContext(ctx)
}

func Debug(args ...interface{}) {
	logging.Print(1, DEBUG, args...)
}

func Debugf(format string, args ...interface{}) {
	logging.Printf(1, DEBUG, format, args...)
}

func Debugw(message string, kvs ...interface{}) {
	logging.Printw(1, DEBUG, message, kvs...)
}

func Info(args ...interface{}) {
	logging.Print(1, INFO, args...)
}

func Infof(format string, args ...interface{}) {
	logging.Printf(1, INFO, format, args...)
}

func Infow(message string, kvs ...interface{}) {
	logging.Printw(1, INFO, message, kvs...)
}

func Warn(args ...interface{}) {
	logging.Print(1, WARN, args...)
}

func Warnf(format string, args ...interface{}) {
	logging.Printf(1, WARN, format, args...)
}

func Warnw(message string, kvs ...interface{}) {
	logging.Printw(1, WARN, message, kvs...)
}

func Error(args ...interface{}) {
	logging.Print(1, ERROR, args...)
}

func Errorf(format string, args ...interface{}) {
	logging.Printf(1, ERROR, format, args...)
}

func Errorw(message string, kvs ...interface{}) {
	logging.Printw(1, ERROR, message, kvs...)
}

func Panic(args ...interface{}) {
	logging.Print(1, PANIC, args...)
}

func Panicf(format string, args ...interface{}) {
	logging.Printf(1, PANIC, format, args...)
}

func Panicw(message string, kvs ...interface{}) {
	logging.Printw(1, PANIC, message, kvs...)
}

func Fatal(args ...interface{}) {
	logging.Print(1, FATAL, args...)
}

func Fatalf(format string, args ...interface{}) {
	logging.Printf(1, FATAL, format, args...)
}

func Fatalw(message string, kvs ...interface{}) {
	logging.Printw(1, FATAL, message, kvs...)
}

func Print(depth int, level Level, args ...interface{}) {
	logging.Print(depth+1, level, args...)
}

func Printf(depth int, level Level, format string, args ...interface{}) {
	logging.Printf(depth+1, level, format, args...)
}

func Printw(depth int, level Level, message string, kvs ...interface{}) {
	logging.Printw(depth+1, level, message, kvs...)
}
