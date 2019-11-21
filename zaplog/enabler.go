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
	if p.min <= lvl && lvl <= p.max {
		return p.level.Enabled(lvl)
	}
	return false
}
