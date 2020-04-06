package zbase

import (
	"testing"

	"go.uber.org/zap/zapcore"

	"git.xiaojukeji.com/pearls/tlog/iface"
)

func TestZapLevel(t *testing.T) {
	tests := []struct {
		ilevel iface.Level
		zlevel zapcore.Level
	}{
		{ilevel: -5, zlevel: zapcore.DebugLevel},
		{ilevel: iface.DEBUG, zlevel: zapcore.DebugLevel},
		{ilevel: iface.INFO, zlevel: zapcore.InfoLevel},
		{ilevel: iface.WARN, zlevel: zapcore.WarnLevel},
		{ilevel: iface.ERROR, zlevel: zapcore.ErrorLevel},
		{ilevel: iface.PANIC, zlevel: zapcore.PanicLevel},
		{ilevel: iface.FATAL, zlevel: zapcore.FatalLevel},
		{ilevel: 10, zlevel: zapcore.FatalLevel},
	}
	for i, tt := range tests {
		level := ZapLevel(tt.ilevel)
		if got, want := level, tt.zlevel; got != want {
			t.Errorf("%d: zap level: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: iface.Level=%v, zapcore.Level=%v", i, tt.ilevel, level)
	}
}

func TestLogLevel(t *testing.T) {
	tests := []struct {
		zlevel zapcore.Level
		ilevel iface.Level
	}{
		{zlevel: -5, ilevel: iface.DEBUG},
		{zlevel: zapcore.DebugLevel, ilevel: iface.DEBUG},
		{zlevel: zapcore.InfoLevel, ilevel: iface.INFO},
		{zlevel: zapcore.WarnLevel, ilevel: iface.WARN},
		{zlevel: zapcore.ErrorLevel, ilevel: iface.ERROR},
		{zlevel: zapcore.DPanicLevel, ilevel: iface.PANIC},
		{zlevel: zapcore.PanicLevel, ilevel: iface.PANIC},
		{zlevel: zapcore.FatalLevel, ilevel: iface.FATAL},
		{zlevel: 10, ilevel: iface.FATAL},
	}
	for i, tt := range tests {
		level := LogLevel(tt.zlevel)
		if got, want := level, tt.ilevel; got != want {
			t.Errorf("%d: iface level: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: zapcore.Level=%v, iface.Level=%v", i, tt.zlevel, level)
	}
}
