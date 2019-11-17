package zaplog

import (
	"errors"
	"fmt"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ironzhang/tlog/iface"
	"github.com/ironzhang/tlog/zaplog/zbase"
	"github.com/ironzhang/tlog/zaplog/zlogger"
)

type Logger struct {
	iface.Logger
	level   zap.AtomicLevel
	sinks   map[string]zap.Sink
	cores   map[string]zapcore.Core
	loggers map[string]iface.Logger
}

func New(cfg Config) (*Logger, error) {
	var logger Logger
	if err := logger.init(cfg); err != nil {
		return nil, err
	}
	return &logger, nil
}

func (p *Logger) init(cfg Config) (err error) {
	p.level = zap.NewAtomicLevelAt(zbase.ZapLevel(cfg.Level))
	for _, sink := range cfg.Sinks {
		if err = p.openSink(sink); err != nil {
			return err
		}
	}
	for _, core := range cfg.Cores {
		if err = p.openCore(core); err != nil {
			return err
		}
	}
	for _, logger := range cfg.Loggers {
		if err = p.openLogger(logger); err != nil {
			return err
		}
	}
	def, ok := p.loggers["default"]
	if !ok {
		return errors.New("not find default logger")
	}
	p.Logger = def
	return nil
}

func (p *Logger) openSink(cfg SinkConfig) error {
	if _, ok := p.sinks[cfg.Name]; ok {
		return fmt.Errorf("sink %q is already opened", cfg.Name)
	}
	sink, err := newSink(cfg.Name, cfg.URL)
	if err != nil {
		return fmt.Errorf("new sink: %w", err)
	}
	p.sinks[cfg.Name] = sink
	return nil
}

func (p *Logger) openCore(cfg CoreConfig) error {
	if _, ok := p.cores[cfg.Name]; ok {
		return fmt.Errorf("core %q is already opened", cfg.Name)
	}

	enc, err := newEncoder(cfg.Encoding, cfg.Encoder)
	if err != nil {
		return fmt.Errorf("new encoder: %w", err)
	}

	ws, err := p.combineSink(cfg.Sinks)
	if err != nil {
		return fmt.Errorf("combine sink: %w", err)
	}

	enab := &levelEnabler{
		min:   zbase.ZapLevel(cfg.MinLevel),
		max:   zbase.ZapLevel(cfg.MaxLevel),
		level: p.level,
	}

	p.cores[cfg.Name] = zapcore.NewCore(enc, ws, enab)

	return nil
}

func (p *Logger) openLogger(cfg LoggerConfig) error {
	if _, ok := p.loggers[cfg.Name]; ok {
		return fmt.Errorf("logger %q is already opened", cfg.Name)
	}

	core, err := p.combineCore(cfg.Cores)
	if err != nil {
		return fmt.Errorf("combine core: %w", err)
	}
	p.loggers[cfg.Name] = zlogger.New(cfg.Name, core, nil)

	return nil
}

func (p *Logger) combineSink(names []string) (zapcore.WriteSyncer, error) {
	sinks := make([]zapcore.WriteSyncer, 0, len(names))
	for _, name := range names {
		sink, ok := p.sinks[name]
		if !ok {
			return nil, fmt.Errorf("not found sink %q", name)
		}
		sinks = append(sinks, sink)
	}
	return zap.CombineWriteSyncers(sinks...), nil
}

func (p *Logger) combineCore(names []string) (zapcore.Core, error) {
	cores := make([]zapcore.Core, 0, len(names))
	for _, name := range names {
		core, ok := p.cores[name]
		if !ok {
			return nil, fmt.Errorf("not found core %q", name)
		}
		cores = append(cores, core)
	}
	return zapcore.NewTee(cores...), nil
}

func (p *Logger) Close() error {
	var err error
	for _, sink := range p.sinks {
		err = multierr.Append(err, sink.Close())
	}
	return err
}

func (p *Logger) Sync() error {
	var err error
	for _, core := range p.cores {
		err = multierr.Append(err, core.Sync())
	}
	return err
}

func (p *Logger) GetLevel() iface.Level {
	return zbase.LogLevel(p.level.Level())
}

func (p *Logger) SetLevel(level iface.Level) {
	p.level.SetLevel(zbase.ZapLevel(level))
}

func (p *Logger) GetDefaultLogger() iface.Logger {
	return p.Logger
}

func (p *Logger) GetLogger(name string) iface.Logger {
	if logger, ok := p.loggers[name]; ok {
		return logger
	}
	return p.Logger
}
