package zlogger

import (
	"go.uber.org/zap/zapcore"

	"github.com/ironzhang/tlog/logger"
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
