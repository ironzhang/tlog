package zaplog

import (
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ironzhang/tlog/logger"
	"github.com/ironzhang/tlog/zaplog/zbase"
	"github.com/ironzhang/tlog/zaplog/zconfig"
)

type ZapLogger struct {
	logger.Logger
	level   zap.AtomicLevel
	sinks   map[string]zap.Sink
	cores   map[string]zapcore.Core
	loggers map[string]logger.Logger
}

func (p *ZapLogger) init(cfg zconfig.Config) error {
	p.level = zap.NewAtomicLevelAt(zbase.ZapLevel(cfg.Level))
	return nil
}

func (p *ZapLogger) Close() error {
	var err error
	for _, sink := range p.sinks {
		err = multierr.Append(err, sink.Close())
	}
	return err
}

func (p *ZapLogger) Sync() error {
	var err error
	for _, core := range p.cores {
		err = multierr.Append(err, core.Sync())
	}
	return err
}

func (p *ZapLogger) GetLevel() logger.Level {
	return zbase.LoggerLevel(p.level.Level())
}

func (p *ZapLogger) SetLevel(level logger.Level) {
	p.level.SetLevel(zbase.ZapLevel(level))
}

func (p *ZapLogger) GetDefaultLogger() logger.Logger {
	return p.Logger
}

func (p *ZapLogger) GetLogger(name string) logger.Logger {
	if logger, ok := p.loggers[name]; ok {
		return logger
	}
	return p.Logger
}
