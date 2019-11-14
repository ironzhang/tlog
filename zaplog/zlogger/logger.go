package zlogger

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ironzhang/tlog/logger"
)

type WithContextFunc func(ctx context.Context) (args []interface{})

type Logger struct {
	name    string
	base    *zap.Logger
	level   zap.AtomicLevel
	withctx WithContextFunc

	ctxs []interface{}
	args []interface{}
}

func New(name string, level Level, core zapcore.Core, opts ...zap.Option) *Logger {
	lv := zap.NewAtomicLevelAt(zapLevel(level))
	base := zap.New(newEnabledCore(core, lv)).WithOptions(opts...)
	return &Logger{
		name:  name,
		base:  base,
		level: lv,
	}
}

func (p *Logger) clone(nctxs, nargs int) *Logger {
	c := &Logger{
		name:    p.name,
		base:    p.base,
		level:   p.level,
		withctx: p.withctx,
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

func (p *Logger) Name() string {
	return p.name
}

func (p *Logger) GetLevel() Level {
	return logLevel(p.level.Level())
}

func (p *Logger) SetLevel(level Level) {
	p.level.SetLevel(zapLevel(level))
}

func (p *Logger) SetWithContextFunc(f WithContextFunc) {
	p.withctx = f
}

func (p *Logger) WithOptions(opts ...zap.Option) *Logger {
	c := p.clone(0, 0)
	c.base = c.base.WithOptions(opts...)
	return c
}

func (p *Logger) WithArgs(args ...interface{}) logger.Logger {
	if len(args) <= 0 {
		return p
	}
	c := p.clone(0, len(args))
	c.args = append(c.args, args...)
	return c
}

func (p *Logger) WithContext(ctx context.Context) logger.Logger {
	if p.withctx == nil {
		return p
	}
	args := p.withctx(ctx)
	if len(args) <= 0 {
		return p
	}
	c := p.clone(len(args), 0)
	c.ctxs = append(c.ctxs, args...)
	return c
}

func (p *Logger) Debug(args ...interface{}) {
	p.Print(1, DEBUG, args...)
}

func (p *Logger) Debugf(format string, args ...interface{}) {
	p.Printf(1, DEBUG, format, args...)
}

func (p *Logger) Debugw(message string, kvs ...interface{}) {
	p.Printw(1, DEBUG, message, kvs...)
}

func (p *Logger) Info(args ...interface{}) {
	p.Print(1, INFO, args...)
}

func (p *Logger) Infof(format string, args ...interface{}) {
	p.Printf(1, INFO, format, args...)
}

func (p *Logger) Infow(message string, kvs ...interface{}) {
	p.Printw(1, INFO, message, kvs...)
}

func (p *Logger) Warn(args ...interface{}) {
	p.Print(1, WARN, args...)
}

func (p *Logger) Warnf(format string, args ...interface{}) {
	p.Printf(1, WARN, format, args...)
}

func (p *Logger) Warnw(message string, kvs ...interface{}) {
	p.Printw(1, WARN, message, kvs...)
}

func (p *Logger) Error(args ...interface{}) {
	p.Print(1, ERROR, args...)
}

func (p *Logger) Errorf(format string, args ...interface{}) {
	p.Printf(1, ERROR, format, args...)
}

func (p *Logger) Errorw(message string, kvs ...interface{}) {
	p.Printw(1, ERROR, message, kvs...)
}

func (p *Logger) Panic(args ...interface{}) {
	p.Print(1, PANIC, args...)
}

func (p *Logger) Panicf(format string, args ...interface{}) {
	p.Printf(1, PANIC, format, args...)
}

func (p *Logger) Panicw(message string, kvs ...interface{}) {
	p.Printw(1, PANIC, message, kvs...)
}

func (p *Logger) Fatal(args ...interface{}) {
	p.Print(1, FATAL, args...)
}

func (p *Logger) Fatalf(format string, args ...interface{}) {
	p.Printf(1, FATAL, format, args...)
}

func (p *Logger) Fatalw(message string, kvs ...interface{}) {
	p.Printw(1, FATAL, message, kvs...)
}

func (p *Logger) Print(depth int, level Level, args ...interface{}) {
	p.log(depth, zapLevel(level), "", args, nil)
}

func (p *Logger) Printf(depth int, level Level, format string, args ...interface{}) {
	p.log(depth, zapLevel(level), format, args, nil)
}

func (p *Logger) Printw(depth int, level Level, message string, kvs ...interface{}) {
	p.log(depth, zapLevel(level), message, nil, kvs)
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
