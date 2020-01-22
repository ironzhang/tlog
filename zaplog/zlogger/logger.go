package zlogger

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ironzhang/tlog/iface"
	"github.com/ironzhang/tlog/zaplog/zbase"
)

type ContextHook interface {
	WithContext(ctx context.Context) (args []interface{})
}

type Logger struct {
	base *zap.Logger
	hook ContextHook

	ctxs []interface{}
	args []interface{}
}

func New(name string, core zapcore.Core, hook ContextHook, opts ...zap.Option) *Logger {
	base := zap.New(core, opts...).Named(name)
	return &Logger{
		base: base,
		hook: hook,
	}
}

func (p *Logger) clone(nctxs, nargs int) *Logger {
	c := &Logger{
		base: p.base,
		hook: p.hook,
	}
	if n := len(p.ctxs); n > 0 {
		c.ctxs = make([]interface{}, n, n+nctxs)
		copy(c.ctxs, p.ctxs)
	}
	if n := len(p.args); n > 0 {
		c.args = make([]interface{}, n, n+nargs)
		copy(c.args, p.args)
	}
	return c
}

func (p *Logger) Named(name string) iface.Logger {
	if len(name) <= 0 {
		return p
	}
	c := p.clone(0, 0)
	c.base = c.base.Named(name)
	return c
}

func (p *Logger) WithArgs(args ...interface{}) iface.Logger {
	if len(args) <= 0 {
		return p
	}
	c := p.clone(0, len(args))
	c.args = append(c.args, args...)
	return c
}

func (p *Logger) WithContext(ctx context.Context) iface.Logger {
	if p.hook == nil {
		return p
	}
	args := p.hook.WithContext(ctx)
	if len(args) <= 0 {
		return p
	}
	c := p.clone(len(args), 0)
	c.ctxs = append(c.ctxs, args...)
	return c
}

func (p *Logger) Sync() error {
	return p.base.Sync()
}

func (p *Logger) Debug(args ...interface{}) {
	p.Print(1, iface.DEBUG, args...)
}

func (p *Logger) Debugf(format string, args ...interface{}) {
	p.Printf(1, iface.DEBUG, format, args...)
}

func (p *Logger) Debugw(message string, kvs ...interface{}) {
	p.Printw(1, iface.DEBUG, message, kvs...)
}

func (p *Logger) Info(args ...interface{}) {
	p.Print(1, iface.INFO, args...)
}

func (p *Logger) Infof(format string, args ...interface{}) {
	p.Printf(1, iface.INFO, format, args...)
}

func (p *Logger) Infow(message string, kvs ...interface{}) {
	p.Printw(1, iface.INFO, message, kvs...)
}

func (p *Logger) Warn(args ...interface{}) {
	p.Print(1, iface.WARN, args...)
}

func (p *Logger) Warnf(format string, args ...interface{}) {
	p.Printf(1, iface.WARN, format, args...)
}

func (p *Logger) Warnw(message string, kvs ...interface{}) {
	p.Printw(1, iface.WARN, message, kvs...)
}

func (p *Logger) Error(args ...interface{}) {
	p.Print(1, iface.ERROR, args...)
}

func (p *Logger) Errorf(format string, args ...interface{}) {
	p.Printf(1, iface.ERROR, format, args...)
}

func (p *Logger) Errorw(message string, kvs ...interface{}) {
	p.Printw(1, iface.ERROR, message, kvs...)
}

func (p *Logger) Panic(args ...interface{}) {
	p.Print(1, iface.PANIC, args...)
}

func (p *Logger) Panicf(format string, args ...interface{}) {
	p.Printf(1, iface.PANIC, format, args...)
}

func (p *Logger) Panicw(message string, kvs ...interface{}) {
	p.Printw(1, iface.PANIC, message, kvs...)
}

func (p *Logger) Fatal(args ...interface{}) {
	p.Print(1, iface.FATAL, args...)
}

func (p *Logger) Fatalf(format string, args ...interface{}) {
	p.Printf(1, iface.FATAL, format, args...)
}

func (p *Logger) Fatalw(message string, kvs ...interface{}) {
	p.Printw(1, iface.FATAL, message, kvs...)
}

func (p *Logger) Print(depth int, level iface.Level, args ...interface{}) {
	p.log(depth, zbase.ZapLevel(level), "", args, nil)
}

func (p *Logger) Printf(depth int, level iface.Level, format string, args ...interface{}) {
	p.log(depth, zbase.ZapLevel(level), format, args, nil)
}

func (p *Logger) Printw(depth int, level iface.Level, message string, kvs ...interface{}) {
	p.log(depth, zbase.ZapLevel(level), message, nil, kvs)
}

func (p *Logger) log(depth int, lvl zapcore.Level, template string, args []interface{}, kvs []interface{}) {
	// If logging at this level is completely disabled, skip the overhead of
	// string formatting.
	if lvl < zapcore.DPanicLevel && !p.base.Core().Enabled(lvl) {
		return
	}

	// Format with Sprint, Sprintf, or neither.
	msg := template
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}

	// Output log message.
	const skip = 2
	base := p.base.WithOptions(zap.AddCallerSkip(skip + depth))
	sugar := base.Sugar().With(p.ctxs...).With(p.args...)
	switch lvl {
	case zapcore.DebugLevel:
		sugar.Debugw(msg, kvs...)
	case zapcore.InfoLevel:
		sugar.Infow(msg, kvs...)
	case zapcore.WarnLevel:
		sugar.Warnw(msg, kvs...)
	case zapcore.ErrorLevel:
		sugar.Errorw(msg, kvs...)
	case zapcore.DPanicLevel:
		sugar.DPanicw(msg, kvs...)
	case zapcore.PanicLevel:
		sugar.Panicw(msg, kvs...)
	case zapcore.FatalLevel:
		sugar.Fatalw(msg, kvs...)
	default:
		if lvl > zapcore.FatalLevel {
			sugar.Fatalw(msg, kvs...)
		} else {
			sugar.Debugw(msg, kvs...)
		}
	}
}
