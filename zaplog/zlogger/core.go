package zlogger

import (
	"go.uber.org/zap/zapcore"
)

type enabledCore struct {
	zapcore.Core
	enabler zapcore.LevelEnabler
}

func newEnabledCore(base zapcore.Core, enabler zapcore.LevelEnabler) zapcore.Core {
	return &enabledCore{
		Core:    base,
		enabler: enabler,
	}
}

func (c *enabledCore) Enabled(lv zapcore.Level) bool {
	return c.enabler.Enabled(lv) && c.Core.Enabled(lv)
}
