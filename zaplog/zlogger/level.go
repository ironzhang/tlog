package zlogger

import (
	"sync/atomic"

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

func zapLevel(lv Level) zapcore.Level {
	switch lv {
	case DEBUG:
		return zapcore.DebugLevel
	case INFO:
		return zapcore.InfoLevel
	case WARN:
		return zapcore.WarnLevel
	case ERROR:
		return zapcore.ErrorLevel
	case PANIC:
		return zapcore.PanicLevel
	case FATAL:
		return zapcore.FatalLevel
	}
	if lv > FATAL {
		return zapcore.FatalLevel
	}
	return zapcore.DebugLevel
}

func logLevel(lv zapcore.Level) Level {
	switch lv {
	case zapcore.DebugLevel:
		return DEBUG
	case zapcore.InfoLevel:
		return INFO
	case zapcore.WarnLevel:
		return WARN
	case zapcore.ErrorLevel:
		return ERROR
	case zapcore.DPanicLevel:
		return PANIC
	case zapcore.PanicLevel:
		return PANIC
	case zapcore.FatalLevel:
		return FATAL
	}
	if lv > zapcore.FatalLevel {
		return FATAL
	}
	return DEBUG
}

type atomicLevel struct {
	lv int32
}

func newAtomicLevel(lv zapcore.Level) *atomicLevel {
	return &atomicLevel{lv: int32(lv)}
}

func (l *atomicLevel) SetLevel(lv zapcore.Level) {
	atomic.StoreInt32(&l.lv, int32(lv))
}

func (l *atomicLevel) GetLevel() zapcore.Level {
	return zapcore.Level(atomic.LoadInt32(&l.lv))
}

func (l *atomicLevel) Enabled(lv zapcore.Level) bool {
	return l.GetLevel() <= lv
}
