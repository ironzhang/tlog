package zbase

import (
	"go.uber.org/zap/zapcore"

	"github.com/ironzhang/tlog/iface"
)

func ZapLevel(lv iface.Level) zapcore.Level {
	switch lv {
	case iface.DEBUG:
		return zapcore.DebugLevel
	case iface.INFO:
		return zapcore.InfoLevel
	case iface.WARN:
		return zapcore.WarnLevel
	case iface.ERROR:
		return zapcore.ErrorLevel
	case iface.PANIC:
		return zapcore.PanicLevel
	case iface.FATAL:
		return zapcore.FatalLevel
	}
	if lv > iface.FATAL {
		return zapcore.FatalLevel
	}
	return zapcore.DebugLevel
}

func LoggerLevel(lv zapcore.Level) iface.Level {
	switch lv {
	case zapcore.DebugLevel:
		return iface.DEBUG
	case zapcore.InfoLevel:
		return iface.INFO
	case zapcore.WarnLevel:
		return iface.WARN
	case zapcore.ErrorLevel:
		return iface.ERROR
	case zapcore.DPanicLevel:
		return iface.PANIC
	case zapcore.PanicLevel:
		return iface.PANIC
	case zapcore.FatalLevel:
		return iface.FATAL
	}
	if lv > zapcore.FatalLevel {
		return iface.FATAL
	}
	return iface.DEBUG
}
