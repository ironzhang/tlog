package zaplog

import (
	"go.uber.org/zap/zapcore"
)

type levelEnabler struct {
	min   zapcore.Level
	max   zapcore.Level
	level zapcore.LevelEnabler
}

func (p *levelEnabler) Enabled(lvl zapcore.Level) bool {
	if lvl < p.min || lvl > p.max {
		return false
	}
	return p.level.Enabled(lvl)
}
