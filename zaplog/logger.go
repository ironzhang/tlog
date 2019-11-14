package zaplog

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/ironzhang/tlog/logger"
	"go.uber.org/zap/zapcore"
)

type Level = logger.Level

const (
	DEBUG = logger.DEBUG
	INFO  = logger.INFO
	WARN  = logger.WARN
	ERROR = logger.ERROR
	PANIC = logger.PANIC
	FATAL = logger.FATAL
)

type WithContextFunc func(ctx context.Context) (args []interface{})

type Logger struct {
	name      string
	core      zapcore.Core
	errout    zapcore.WriteSyncer
	addstack  zapcore.LevelEnabler
	addcaller bool

	withctx WithContextFunc
	ctxs    []zapcore.Field
	args    []zapcore.Field
}

func (p *Logger) GetLevel() Level {
	return DEBUG
}

func (p *Logger) SetLevel(level Level) {
}

func (p *Logger) WithArgs(args ...interface{}) logger.Logger {
	return p
}

func (p *Logger) WithContext(ctx context.Context) logger.Logger {
	return p
}

func (p *Logger) Debug(args ...interface{}) {
}

func (p *Logger) Debugf(format string, args ...interface{}) {
}

func (p *Logger) Debugw(message string, kvs ...interface{}) {
}

func (p *Logger) Info(args ...interface{}) {
}

func (p *Logger) Infof(format string, args ...interface{}) {
}

func (p *Logger) Infow(message string, kvs ...interface{}) {
}

func (p *Logger) Warn(args ...interface{}) {
}

func (p *Logger) Warnf(format string, args ...interface{}) {
}

func (p *Logger) Warnw(message string, kvs ...interface{}) {
}

func (p *Logger) Error(args ...interface{}) {
}

func (p *Logger) Errorf(format string, args ...interface{}) {
}

func (p *Logger) Errorw(message string, kvs ...interface{}) {
}

func (p *Logger) Panic(args ...interface{}) {
}

func (p *Logger) Panicf(format string, args ...interface{}) {
}

func (p *Logger) Panicw(message string, kvs ...interface{}) {
}

func (p *Logger) Fatal(args ...interface{}) {
}

func (p *Logger) Fatalf(format string, args ...interface{}) {
}

func (p *Logger) Fatalw(message string, kvs ...interface{}) {
}

func (p *Logger) Print(depth int, level Level, args ...interface{}) {
}

func (p *Logger) Printf(depth int, level Level, format string, args ...interface{}) {
}

func (p *Logger) Printw(depth int, level Level, message string, kvs ...interface{}) {
}

func (p *Logger) log(depth int, level zapcore.Level, template string, args []interface{}, kvs []interface{}) {
	if level < zapcore.DPanicLevel && !p.core.Enabled(level) {
		return
	}

	msg := template
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}

	core := p.core
	if len(p.ctxs) > 0 {
		core = core.With(p.ctxs)
	}
	if len(p.args) > 0 {
		core = core.With(p.args)
	}

	if ce := p.check(core, depth, level, msg); ce != nil {
		ce.Write(p.sweetenFields(kvs)...)
	}
}

func (p *Logger) check(core zapcore.Core, depth int, level zapcore.Level, msg string) *zapcore.CheckedEntry {
	const callerskip = 2

	ent := zapcore.Entry{
		Level:      level,
		Time:       time.Now(),
		LoggerName: p.name,
		Message:    msg,
	}
	ce := core.Check(ent, nil)
	if ce != nil {
		ce.ErrorOutput = p.errout
		if p.addcaller {
			ce.Entry.Caller = zapcore.NewEntryCaller(runtime.Caller(callerskip + depth))
			if !ce.Entry.Caller.Defined {
				fmt.Fprintf(p.errout, "%v Logger.check error: failed to get caller\n", time.Now().UTC())
				p.errout.Sync()
			}
		}
		if p.addstack.Enabled(ce.Entry.Level) {
			ce.Entry.Stack = stackTrace(callerskip + depth)
		}
	}

	switch ent.Level {
	case zapcore.PanicLevel:
		ce = ce.Should(ent, zapcore.WriteThenPanic)
	case zapcore.FatalLevel:
		ce = ce.Should(ent, zapcore.WriteThenFatal)
	}
	return ce
}

func (p *Logger) sweetenFields(args []interface{}) []zapcore.Field {
	return nil
}
