package zlogger

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ironzhang/tlog/logger"
)

type WithContextFunc func(ctx context.Context) (args []interface{})

type Logger struct {
	name        string
	level       zap.AtomicLevel
	core        zapcore.Core
	errout      zapcore.WriteSyncer
	addstack    zapcore.LevelEnabler
	addcaller   bool
	withContext WithContextFunc

	ctxs []zapcore.Field
	args []zapcore.Field
}

func NewLogger(name string, level Level, core zapcore.Core) *Logger {
	lv := zap.NewAtomicLevelAt(zapLevel(level))
	return &Logger{
		name:      name,
		level:     lv,
		core:      newEnabledCore(core, lv),
		errout:    zapcore.Lock(os.Stdout),
		addstack:  zapcore.ErrorLevel,
		addcaller: true,
	}
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

func (p *Logger) clone() *Logger {
	c := &Logger{
		name:      p.name,
		level:     p.level,
		core:      p.core,
		errout:    p.errout,
		addstack:  p.addstack,
		addcaller: p.addcaller,
	}
	if n := len(p.ctxs); n > 0 {
		c.ctxs = make([]zapcore.Field, n, n+10)
		copy(c.ctxs, p.ctxs)
	}
	if n := len(p.args); n > 0 {
		c.args = make([]zapcore.Field, n, n+10)
		copy(c.args, p.args)
	}
	return c
}

func (p *Logger) WithArgs(args ...interface{}) logger.Logger {
	if len(args) <= 0 {
		return p
	}
	c := p.clone()
	c.args = append(c.args, p.sweetenFields(args)...)
	return c
}

func (p *Logger) WithContext(ctx context.Context) logger.Logger {
	if p.withContext == nil {
		return p
	}
	args := p.withContext(ctx)
	if len(args) <= 0 {
		return p
	}
	c := p.clone()
	c.ctxs = append(c.ctxs, p.sweetenFields(args)...)
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

	p.output(depth, level, msg, p.sweetenFields(kvs)...)
}

func (p *Logger) output(depth int, level zapcore.Level, msg string, fields ...zapcore.Field) {
	core := p.core
	if len(p.ctxs) > 0 {
		core = core.With(p.ctxs)
	}
	if len(p.args) > 0 {
		core = core.With(p.args)
	}
	if ce := p.check(core, depth, level, msg); ce != nil {
		ce.Write(fields...)
	}
}

func (p *Logger) check(core zapcore.Core, depth int, level zapcore.Level, msg string) *zapcore.CheckedEntry {
	const callerskip = 3

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

const (
	_oddNumberErrMsg    = "Ignored key without a value."
	_nonStringKeyErrMsg = "Ignored key-value pairs with non-string keys."
)

func (p *Logger) sweetenFields(args []interface{}) []zapcore.Field {
	if len(args) == 0 {
		return nil
	}

	// Allocate enough space for the worst case; if users pass only structured
	// fields, we shouldn't penalize them with extra allocations.
	fields := make([]zapcore.Field, 0, len(args))
	var invalid invalidPairs

	for i := 0; i < len(args); {
		// This is a strongly-typed field. Consume it and move on.
		if f, ok := args[i].(zapcore.Field); ok {
			fields = append(fields, f)
			i++
			continue
		}

		// Make sure this element isn't a dangling key.
		if i == len(args)-1 {
			p.output(1, zapcore.PanicLevel, _oddNumberErrMsg, zap.Any("ignored", args[i]))
			break
		}

		// Consume this value and the next, treating them as a key-value pair. If the
		// key isn't a string, add this pair to the slice of invalid pairs.
		key, val := args[i], args[i+1]
		if keyStr, ok := key.(string); !ok {
			// Subsequent errors are likely, so allocate once up front.
			if cap(invalid) == 0 {
				invalid = make(invalidPairs, 0, len(args)/2)
			}
			invalid = append(invalid, invalidPair{i, key, val})
		} else {
			fields = append(fields, zap.Any(keyStr, val))
		}
		i += 2
	}

	// If we encountered any invalid key-value pairs, log an error.
	if len(invalid) > 0 {
		p.output(1, zapcore.PanicLevel, _nonStringKeyErrMsg, zap.Array("invalid", invalid))
	}
	return fields
}

type invalidPair struct {
	position   int
	key, value interface{}
}

func (p invalidPair) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("position", int64(p.position))
	zap.Any("key", p.key).AddTo(enc)
	zap.Any("value", p.value).AddTo(enc)
	return nil
}

type invalidPairs []invalidPair

func (ps invalidPairs) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	var err error
	for i := range ps {
		err = multierr.Append(err, enc.AppendObject(ps[i]))
	}
	return err
}
