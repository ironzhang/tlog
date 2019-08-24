package zaplog

import (
	"context"

	"github.com/ironzhang/tlog/logger"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	base *zap.Logger
	hook logger.ContextHookFunc
	ctxs []zap.Field
	args []zap.Field
}

func NewLogger(base *zap.Logger) *Logger {
	return &Logger{
		base: base.WithOptions(zap.AddCallerSkip(1)),
	}
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
	if p.hook == nil {
		return p
	}
	args := p.hook(ctx)
	if len(args) <= 0 {
		return p
	}
	c := p.clone()
	c.ctxs = append(c.ctxs, p.sweetenFields(args)...)
	return c
}

func (p *Logger) Debug(args ...interface{}) {
	p.Sugar().Debug(args...)
}

func (p *Logger) Debugf(format string, args ...interface{}) {
	p.Sugar().Debugf(format, args...)
}

func (p *Logger) Debugw(message string, kvs ...interface{}) {
	p.Sugar().Debugw(message, kvs...)
}

func (p *Logger) Trace(args ...interface{}) {
	p.Sugar().Debug(args...)
}

func (p *Logger) Tracef(format string, args ...interface{}) {
	p.Sugar().Debugf(format, args...)
}

func (p *Logger) Tracew(message string, kvs ...interface{}) {
	p.Sugar().Debugw(message, kvs...)
}

func (p *Logger) Info(args ...interface{}) {
	p.Sugar().Info(args...)
}

func (p *Logger) Infof(format string, args ...interface{}) {
	p.Sugar().Infof(format, args...)
}

func (p *Logger) Infow(message string, kvs ...interface{}) {
	p.Sugar().Infow(message, kvs...)
}

func (p *Logger) Warn(args ...interface{}) {
	p.Sugar().Warn(args...)
}

func (p *Logger) Warnf(format string, args ...interface{}) {
	p.Sugar().Warnf(format, args...)
}

func (p *Logger) Warnw(message string, kvs ...interface{}) {
	p.Sugar().Warnw(message, kvs...)
}

func (p *Logger) Error(args ...interface{}) {
	p.Sugar().Error(args...)
}

func (p *Logger) Errorf(format string, args ...interface{}) {
	p.Sugar().Errorf(format, args...)
}

func (p *Logger) Errorw(message string, kvs ...interface{}) {
	p.Sugar().Errorw(message, kvs...)
}

func (p *Logger) Panic(args ...interface{}) {
	p.Sugar().Panic(args...)
}

func (p *Logger) Panicf(format string, args ...interface{}) {
	p.Sugar().Panicf(format, args...)
}

func (p *Logger) Panicw(message string, kvs ...interface{}) {
	p.Sugar().Panicw(message, kvs...)
}

func (p *Logger) Fatal(args ...interface{}) {
	p.Sugar().Fatal(args...)
}

func (p *Logger) Fatalf(format string, args ...interface{}) {
	p.Sugar().Fatalf(format, args...)
}

func (p *Logger) Fatalw(message string, kvs ...interface{}) {
	p.Sugar().Fatalw(message, kvs...)
}

func (p *Logger) Base() *zap.Logger {
	return p.base.With(p.ctxs...).With(p.args...)
}

func (p *Logger) Sugar() *zap.SugaredLogger {
	return p.Base().Sugar()
}

func (p *Logger) clone() *Logger {
	c := &Logger{
		base: p.base,
		hook: p.hook,
	}
	if len(p.ctxs) > 0 {
		c.ctxs = make([]zap.Field, len(p.ctxs))
		copy(c.ctxs, p.ctxs)
	}
	if len(p.args) > 0 {
		c.args = make([]zap.Field, len(p.args))
		copy(c.args, p.args)
	}
	return c
}

const (
	_oddNumberErrMsg    = "Ignored key without a value."
	_nonStringKeyErrMsg = "Ignored key-value pairs with non-string keys."
)

func (p *Logger) sweetenFields(args []interface{}) []zap.Field {
	if len(args) == 0 {
		return nil
	}

	// Allocate enough space for the worst case; if users pass only structured
	// fields, we shouldn't penalize them with extra allocations.
	fields := make([]zap.Field, 0, len(args))
	var invalid invalidPairs

	for i := 0; i < len(args); {
		// This is a strongly-typed field. Consume it and move on.
		if f, ok := args[i].(zap.Field); ok {
			fields = append(fields, f)
			i++
			continue
		}

		// Make sure this element isn't a dangling key.
		if i == len(args)-1 {
			p.Base().DPanic(_oddNumberErrMsg, zap.Any("ignored", args[i]))
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
		p.Base().DPanic(_nonStringKeyErrMsg, zap.Array("invalid", invalid))
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
