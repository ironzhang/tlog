package zaplog

/*
type enabledCore struct {
	zapcore.Core
	enabler zapcore.LevelEnabler
}

func (c *enabledCore) Enabled(lv zapcore.Level) bool {
	return c.enabler.Enabled(lv) && c.Core.Enabled(lv)
}

func newEnabledCore(base zapcore.Core, enabler zapcore.LevelEnabler) zapcore.Core {
	return &enabledCore{
		Core:    base,
		enabler: enabler,
	}
}
*/

/*
type core struct {
}

func newCore(name string, encoding string, encoder zapcore.EncoderConfig,
	minLevel, maxLevel logger.Level, ws zapcore.WriteSyncer) (*core, error) {
	enc, err := newEncoder(encoding, encoder)
	if err != nil {
		return nil, err
	}
	zapcore.NewCore(enc, ws)
	return nil, nil
}
*/
