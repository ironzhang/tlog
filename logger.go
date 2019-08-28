package tlog

import (
	"context"
)

// ContextHookFunc 上下文钩子函数
type ContextHookFunc func(ctx context.Context) []interface{}

// Logger 日志接口
type Logger interface {
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

type nopLogger struct{}

func (p *nopLogger) WithArgs(args ...interface{}) Logger       { return p }
func (p *nopLogger) WithContext(ctx context.Context) Logger    { return p }
func (p *nopLogger) Trace(args ...interface{})                 {}
func (p *nopLogger) Tracef(format string, args ...interface{}) {}
func (p *nopLogger) Tracew(message string, kvs ...interface{}) {}
func (p *nopLogger) Debug(args ...interface{})                 {}
func (p *nopLogger) Debugf(format string, args ...interface{}) {}
func (p *nopLogger) Debugw(message string, kvs ...interface{}) {}
func (p *nopLogger) Info(args ...interface{})                  {}
func (p *nopLogger) Infof(format string, args ...interface{})  {}
func (p *nopLogger) Infow(message string, kvs ...interface{})  {}
func (p *nopLogger) Warn(args ...interface{})                  {}
func (p *nopLogger) Warnf(format string, args ...interface{})  {}
func (p *nopLogger) Warnw(message string, kvs ...interface{})  {}
func (p *nopLogger) Error(args ...interface{})                 {}
func (p *nopLogger) Errorf(format string, args ...interface{}) {}
func (p *nopLogger) Errorw(message string, kvs ...interface{}) {}
func (p *nopLogger) Panic(args ...interface{})                 {}
func (p *nopLogger) Panicf(format string, args ...interface{}) {}
func (p *nopLogger) Panicw(message string, kvs ...interface{}) {}
func (p *nopLogger) Fatal(args ...interface{})                 {}
func (p *nopLogger) Fatalf(format string, args ...interface{}) {}
func (p *nopLogger) Fatalw(message string, kvs ...interface{}) {}
