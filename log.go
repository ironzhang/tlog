package tlog

import (
	"context"

	"github.com/ironzhang/tlog/logger"
)

var logging logger.Logger = nopLogger{}

var (
	Named       = logging.Named
	WithArgs    = logging.WithArgs
	WithContext = logging.WithContext

	Trace  = logging.Trace
	Tracef = logging.Tracef
	Tracew = logging.Tracew

	Debug  = logging.Debug
	Debugf = logging.Debugf
	Debugw = logging.Debugw

	Info  = logging.Info
	Infof = logging.Infof
	Infow = logging.Infow

	Warn  = logging.Warn
	Warnf = logging.Warnf
	Warnw = logging.Warnw

	Error  = logging.Error
	Errorf = logging.Errorf
	Errorw = logging.Errorw

	Panic  = logging.Panic
	Panicf = logging.Panicf
	Panicw = logging.Panicw

	Fatal  = logging.Fatal
	Fatalf = logging.Fatalf
	Fatalw = logging.Fatalw
)

func SetLogger(l logger.Logger) logger.Logger {
	prev := logging
	if l == nil {
		logging = nopLogger{}
	} else {
		logging = l
	}

	Named = logging.Named
	WithArgs = logging.WithArgs
	WithContext = logging.WithContext
	Trace = logging.Trace
	Tracef = logging.Tracef
	Tracew = logging.Tracew
	Debug = logging.Debug
	Debugf = logging.Debugf
	Debugw = logging.Debugw
	Info = logging.Info
	Infof = logging.Infof
	Infow = logging.Infow
	Warn = logging.Warn
	Warnf = logging.Warnf
	Warnw = logging.Warnw
	Error = logging.Error
	Errorf = logging.Errorf
	Errorw = logging.Errorw
	Panic = logging.Panic
	Panicf = logging.Panicf
	Panicw = logging.Panicw
	Fatal = logging.Fatal
	Fatalf = logging.Fatalf
	Fatalw = logging.Fatalw

	return prev
}

type nopLogger struct{}

func (p nopLogger) Named(name string) logger.Logger               { return p }
func (p nopLogger) WithArgs(args ...interface{}) logger.Logger    { return p }
func (p nopLogger) WithContext(ctx context.Context) logger.Logger { return p }
func (p nopLogger) Trace(args ...interface{})                     {}
func (p nopLogger) Tracef(format string, args ...interface{})     {}
func (p nopLogger) Tracew(message string, kvs ...interface{})     {}
func (p nopLogger) Debug(args ...interface{})                     {}
func (p nopLogger) Debugf(format string, args ...interface{})     {}
func (p nopLogger) Debugw(message string, kvs ...interface{})     {}
func (p nopLogger) Info(args ...interface{})                      {}
func (p nopLogger) Infof(format string, args ...interface{})      {}
func (p nopLogger) Infow(message string, kvs ...interface{})      {}
func (p nopLogger) Warn(args ...interface{})                      {}
func (p nopLogger) Warnf(format string, args ...interface{})      {}
func (p nopLogger) Warnw(message string, kvs ...interface{})      {}
func (p nopLogger) Error(args ...interface{})                     {}
func (p nopLogger) Errorf(format string, args ...interface{})     {}
func (p nopLogger) Errorw(message string, kvs ...interface{})     {}
func (p nopLogger) Panic(args ...interface{})                     {}
func (p nopLogger) Panicf(format string, args ...interface{})     {}
func (p nopLogger) Panicw(message string, kvs ...interface{})     {}
func (p nopLogger) Fatal(args ...interface{})                     {}
func (p nopLogger) Fatalf(format string, args ...interface{})     {}
func (p nopLogger) Fatalw(message string, kvs ...interface{})     {}
