package zlogger

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"git.xiaojukeji.com/pearls/tlog/iface"
	"git.xiaojukeji.com/pearls/tlog/zaplog/zbase"
)

type TContextHook struct {
	trace bool
}

func (p *TContextHook) WithContext(ctx context.Context) (args []interface{}) {
	if p.trace {
		args = append(args, "trace_id", "123456")
	}
	return args
}

type TLogged struct {
	entries []observer.LoggedEntry
}

func (p *TLogged) Add(e zapcore.Entry, fields ...zapcore.Field) {
	if len(fields) == 0 {
		fields = []zapcore.Field{}
	}
	p.entries = append(p.entries, observer.LoggedEntry{Entry: e, Context: fields})
}

func NewTestLogger(t testing.TB, name string, level iface.Level, hook ContextHook,
	opts ...zap.Option) (*Logger, *observer.ObservedLogs) {
	core, logs := observer.New(zbase.ZapLevel(level))
	return New(name, core, hook, opts...), logs
}

func TestLoggerPrint(t *testing.T) {
	var logged TLogged
	logger, logs := NewTestLogger(t, "", iface.DEBUG, nil)

	min, max := iface.DEBUG-1, iface.ERROR
	for lvl := min; lvl <= max; lvl++ {
		logger.Print(0, lvl, "Print", lvl)
		logger.Printf(0, lvl, "Print: %v", lvl)
		logger.Printw(0, lvl, "Print", "level", lvl)

		logged.Add(zapcore.Entry{Level: zbase.ZapLevel(lvl), Message: fmt.Sprint("Print", lvl)})
		logged.Add(zapcore.Entry{Level: zbase.ZapLevel(lvl), Message: fmt.Sprintf("Print: %v", lvl)})
		logged.Add(zapcore.Entry{Level: zbase.ZapLevel(lvl), Message: "Print"}, zap.Any("level", lvl))
	}

	assert.Equal(t, logged.entries, logs.AllUntimed(), "unexpected log entries")
}

func TestLoggerLevel(t *testing.T) {
	var logged TLogged
	logger, logs := NewTestLogger(t, "", iface.WARN, nil)

	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")

	logged.Add(zapcore.Entry{Level: zapcore.WarnLevel, Message: "warn"})
	logged.Add(zapcore.Entry{Level: zapcore.ErrorLevel, Message: "error"})

	assert.Equal(t, logged.entries, logs.AllUntimed(), "unexpected log entries")
}

func TestLoggerLog(t *testing.T) {
	var logged TLogged
	logger, logs := NewTestLogger(t, "", iface.DEBUG, nil)

	tests := []struct {
		level iface.Level
		log   func(args ...interface{})
		logf  func(format string, args ...interface{})
		logw  func(message string, kvs ...interface{})
	}{
		{
			level: iface.DEBUG,
			log:   logger.Debug,
			logf:  logger.Debugf,
			logw:  logger.Debugw,
		},
		{
			level: iface.INFO,
			log:   logger.Info,
			logf:  logger.Infof,
			logw:  logger.Infow,
		},
		{
			level: iface.WARN,
			log:   logger.Warn,
			logf:  logger.Warnf,
			logw:  logger.Warnw,
		},
		{
			level: iface.ERROR,
			log:   logger.Error,
			logf:  logger.Errorf,
			logw:  logger.Errorw,
		},
	}
	for _, tt := range tests {
		tt.log("hello", "world")
		logged.Add(zapcore.Entry{Level: zbase.ZapLevel(tt.level), Message: "helloworld"})
		tt.logf("%s, %s", "hello", "world")
		logged.Add(zapcore.Entry{Level: zbase.ZapLevel(tt.level), Message: "hello, world"})
		tt.logw("greeting", "hello", "world")
		logged.Add(zapcore.Entry{Level: zbase.ZapLevel(tt.level), Message: "greeting"}, zap.String("hello", "world"))
	}

	assert.Equal(t, logged.entries, logs.AllUntimed(), "unexpected log entries")
}

func TestLoggerPanic(t *testing.T) {
	var logged TLogged
	logger, logs := NewTestLogger(t, "", iface.DEBUG, nil)

	assert.Panics(t, func() { logger.Panic("panic") }, "Expected panic")
	assert.Panics(t, func() { logger.Panicf("%s", "panicf") }, "Expected panicf")
	assert.Panics(t, func() { logger.Panicw("panicw", "hello", "world") }, "Expected panicw")

	logged.Add(zapcore.Entry{Level: zapcore.PanicLevel, Message: "panic"})
	logged.Add(zapcore.Entry{Level: zapcore.PanicLevel, Message: "panicf"})
	logged.Add(zapcore.Entry{Level: zapcore.PanicLevel, Message: "panicw"}, zap.String("hello", "world"))

	assert.Equal(t, logged.entries, logs.AllUntimed(), "unexpected log entries")
}

func TestLoggerWithArgs(t *testing.T) {
	var logged TLogged
	logger, logs := NewTestLogger(t, "", iface.DEBUG, nil)

	logger.WithArgs().Info()
	logger.WithArgs("k1", "v1").Info()
	logger.WithArgs("k2", "v2").Info()
	logger.WithArgs("k3", "v3").WithArgs("k4", "v4").Info()

	logged.Add(zapcore.Entry{Level: zapcore.InfoLevel})
	logged.Add(zapcore.Entry{Level: zapcore.InfoLevel}, zap.String("k1", "v1"))
	logged.Add(zapcore.Entry{Level: zapcore.InfoLevel}, zap.String("k2", "v2"))
	logged.Add(zapcore.Entry{Level: zapcore.InfoLevel}, zap.String("k3", "v3"), zap.String("k4", "v4"))

	assert.Equal(t, logged.entries, logs.AllUntimed(), "unexpected log entries")
}

func TestLoggerWithoutContext(t *testing.T) {
	var logged TLogged
	logger, logs := NewTestLogger(t, "", iface.DEBUG, nil)

	logger.WithContext(context.Background()).Info()
	logger.WithContext(context.Background()).WithArgs("k1", "v1").Info()

	logged.Add(zapcore.Entry{Level: zapcore.InfoLevel})
	logged.Add(zapcore.Entry{Level: zapcore.InfoLevel}, zap.String("k1", "v1"))

	assert.Equal(t, logged.entries, logs.AllUntimed(), "unexpected log entries")
}

func TestLoggerWithContext(t *testing.T) {
	hook := TContextHook{}

	var logged TLogged
	logger, logs := NewTestLogger(t, "", iface.DEBUG, &hook)

	logger.WithContext(context.Background()).Info()
	hook.trace = true
	logger.WithContext(context.Background()).Info()
	logger.WithContext(context.Background()).WithArgs("k1", "v1").Info()
	logger.WithArgs("k2", "v2").WithContext(context.Background()).Info()

	logged.Add(zapcore.Entry{Level: zapcore.InfoLevel})
	logged.Add(zapcore.Entry{Level: zapcore.InfoLevel}, zap.String("trace_id", "123456"))
	logged.Add(zapcore.Entry{Level: zapcore.InfoLevel}, zap.String("trace_id", "123456"), zap.String("k1", "v1"))
	logged.Add(zapcore.Entry{Level: zapcore.InfoLevel}, zap.String("trace_id", "123456"), zap.String("k2", "v2"))

	assert.Equal(t, logged.entries, logs.AllUntimed(), "unexpected log entries")
}

func TestLoggerName(t *testing.T) {
	var logged TLogged
	logger, logs := NewTestLogger(t, "name", iface.DEBUG, nil)

	logger.Info()
	logged.Add(zapcore.Entry{Level: zapcore.InfoLevel, LoggerName: "name"})
	assert.Equal(t, logged.entries, logs.AllUntimed(), "unexpected log entries")
}

func TestLoggerAddCaller(t *testing.T) {
	logger, logs := NewTestLogger(t, "name", iface.DEBUG, nil, zap.AddCaller())

	logger.Info()
	output := logs.AllUntimed()
	assert.Equal(t, 1, len(output), "unexpected number of logs written out")
	assert.Regexp(t, `logger_test.go`, output[0].Caller.String(), "unexpected caller")
	//assert.Equal(t, 164, output[0].Caller.Line, "unexpected line")
}
