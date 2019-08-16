package tlog

import (
	"context"
	"log"
	"os"

	"github.com/ironzhang/tlog/logger"
	"github.com/ironzhang/tlog/stdlog"
)

var Default *stdlog.Logger
var logging logger.Logger

func init() {
	Default = stdlog.NewLogger(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile), stdlog.SetCalldepth(1))
	logging = Default
}

func SetLogger(l logger.Logger) logger.Logger {
	prev := logging
	if l == nil {
		logging = Default
	} else {
		logging = l
	}
	return prev
}

func WithArgs(args ...interface{}) logger.Logger {
	return logging.WithArgs(args...)
}

func WithContext(ctx context.Context) logger.Logger {
	return logging.WithContext(ctx)
}

func Debug(args ...interface{}) {
	logging.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logging.Debugf(format, args...)
}

func Debugw(message string, kvs ...interface{}) {
	logging.Debugw(message, kvs...)
}

func Trace(args ...interface{}) {
	logging.Trace(args...)
}

func Tracef(format string, args ...interface{}) {
	logging.Tracef(format, args...)
}

func Tracew(message string, kvs ...interface{}) {
	logging.Tracew(message, kvs...)
}

func Info(args ...interface{}) {
	logging.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logging.Infof(format, args...)
}

func Infow(message string, kvs ...interface{}) {
	logging.Infow(message, kvs...)
}

func Warn(args ...interface{}) {
	logging.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	logging.Warnf(format, args...)
}

func Warnw(message string, kvs ...interface{}) {
	logging.Warnw(message, kvs...)
}

func Error(args ...interface{}) {
	logging.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logging.Errorf(format, args...)
}

func Errorw(message string, kvs ...interface{}) {
	logging.Errorw(message, kvs...)
}

func Panic(args ...interface{}) {
	logging.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	logging.Panicf(format, args...)
}

func Panicw(message string, kvs ...interface{}) {
	logging.Panicw(message, kvs...)
}

func Fatal(args ...interface{}) {
	logging.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	logging.Fatalf(format, args...)
}

func Fatalw(message string, kvs ...interface{}) {
	logging.Fatalw(message, kvs...)
}
