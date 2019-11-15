package zbase

import (
	"go.uber.org/zap/zapcore"

	"github.com/ironzhang/tlog/logger"
)

func ZapLevel(lv logger.Level) zapcore.Level {
	switch lv {
	case logger.DEBUG:
		return zapcore.DebugLevel
	case logger.INFO:
		return zapcore.InfoLevel
	case logger.WARN:
		return zapcore.WarnLevel
	case logger.ERROR:
		return zapcore.ErrorLevel
	case logger.PANIC:
		return zapcore.PanicLevel
	case logger.FATAL:
		return zapcore.FatalLevel
	}
	if lv > logger.FATAL {
		return zapcore.FatalLevel
	}
	return zapcore.DebugLevel
}

func LoggerLevel(lv zapcore.Level) logger.Level {
	switch lv {
	case zapcore.DebugLevel:
		return logger.DEBUG
	case zapcore.InfoLevel:
		return logger.INFO
	case zapcore.WarnLevel:
		return logger.WARN
	case zapcore.ErrorLevel:
		return logger.ERROR
	case zapcore.DPanicLevel:
		return logger.PANIC
	case zapcore.PanicLevel:
		return logger.PANIC
	case zapcore.FatalLevel:
		return logger.FATAL
	}
	if lv > zapcore.FatalLevel {
		return logger.FATAL
	}
	return logger.DEBUG
}
