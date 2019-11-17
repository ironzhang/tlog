package zaplog

import (
	"fmt"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ironzhang/tlog/logger"
	"github.com/ironzhang/tlog/zaplog/zbase"
	"github.com/ironzhang/tlog/zaplog/zlogger"
)

type ZapLogger struct {
	logger.Logger
	level   zap.AtomicLevel
	sinks   map[string]zap.Sink
	cores   map[string]zapcore.Core
	loggers map[string]logger.Logger
}

func (p *ZapLogger) init(cfg Config) (err error) {
	p.level = zap.NewAtomicLevelAt(zbase.ZapLevel(cfg.Level))
	for _, sink := range cfg.Sinks {
		if err = p.addSink(sink); err != nil {
			return err
		}
	}
	for _, core := range cfg.Cores {
		if err = p.addCore(core); err != nil {
			return err
		}
	}
	for _, logger := range cfg.Loggers {
		if err = p.addLogger(logger); err != nil {
			return err
		}
	}
	return nil
}

func (p *ZapLogger) addSink(cfg SinkConfig) error {
	if _, ok := p.sinks[cfg.Name]; ok {
		return fmt.Errorf("sink %q is opened", cfg.Name)
	}
	sink, err := newSink(cfg.Name, cfg.URL)
	if err != nil {
		return err
	}
	p.sinks[cfg.Name] = sink
	return nil
}

func (p *ZapLogger) addCore(cfg CoreConfig) error {
	if _, ok := p.cores[cfg.Name]; ok {
		return fmt.Errorf("core %q is opened", cfg.Name)
	}

	enc, err := newEncoder(cfg.Encoding, cfg.Encoder)
	if err != nil {
		return err
	}

	ws, err := p.combineSink(cfg.Sinks)
	if err != nil {
		return err
	}

	enab := &levelEnabler{
		min:   zbase.ZapLevel(cfg.MinLevel),
		max:   zbase.ZapLevel(cfg.MaxLevel),
		level: p.level,
	}

	p.cores[cfg.Name] = zapcore.NewCore(enc, ws, enab)

	return nil
}

func (p *ZapLogger) addLogger(cfg LoggerConfig) error {
	if _, ok := p.loggers[cfg.Name]; ok {
		return fmt.Errorf("logger %q is opened", cfg.Name)
	}

	core, err := p.combineCore(cfg.Cores)
	if err != nil {
		return err
	}
	p.loggers[cfg.Name] = zlogger.New(cfg.Name, core, nil)

	return nil
}

func (p *ZapLogger) combineSink(names []string) (zapcore.WriteSyncer, error) {
	var sinks []zapcore.WriteSyncer
	for _, name := range names {
		sink, ok := p.sinks[name]
		if !ok {
			return nil, fmt.Errorf("not found sink %q", name)
		}
		sinks = append(sinks, sink)
	}
	return zap.CombineWriteSyncers(sinks...), nil
}

func (p *ZapLogger) combineCore(names []string) (zapcore.Core, error) {
	var cores []zapcore.Core
	for _, name := range names {
		core, ok := p.cores[name]
		if !ok {
			return nil, fmt.Errorf("not found core %q", name)
		}
		cores = append(cores, core)
	}
	return zapcore.NewTee(cores...), nil
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
