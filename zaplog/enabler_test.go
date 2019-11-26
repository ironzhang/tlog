package zaplog

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestLevelEnabler(t *testing.T) {
	tests := []struct {
		enab    levelEnabler
		level   zapcore.Level
		enabled bool
	}{
		{
			enab:    levelEnabler{min: zapcore.DebugLevel, max: zapcore.WarnLevel, level: zapcore.DebugLevel},
			level:   zapcore.DebugLevel,
			enabled: true,
		},
		{
			enab:    levelEnabler{min: zapcore.InfoLevel, max: zapcore.WarnLevel, level: zapcore.DebugLevel},
			level:   zapcore.DebugLevel,
			enabled: false,
		},
		{
			enab:    levelEnabler{min: zapcore.DebugLevel, max: zapcore.WarnLevel, level: zapcore.DebugLevel},
			level:   zapcore.FatalLevel,
			enabled: false,
		},
		{
			enab:    levelEnabler{min: zapcore.DebugLevel, max: zapcore.WarnLevel, level: zapcore.InfoLevel},
			level:   zapcore.DebugLevel,
			enabled: false,
		},
	}
	for i, tt := range tests {
		enabled := tt.enab.Enabled(tt.level)
		if got, want := enabled, tt.enabled; got != want {
			t.Errorf("%d: enabled: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: enabled: got %v", i, enabled)
	}
}
