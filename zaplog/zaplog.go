package zaplog

import (
	"context"
	"errors"
	"fmt"
	"io"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ironzhang/tlog/iface"
	"github.com/ironzhang/tlog/zaplog/zbase"
	"github.com/ironzhang/tlog/zaplog/zlogger"
)

type ContextHook = zlogger.ContextHook

type ContextHookFunc func(ctx context.Context) (args []interface{})

func (f ContextHookFunc) WithContext(ctx context.Context) (args []interface{}) {
	return f(ctx)
}

type Logger struct {
	hook  ContextHook
	level zap.AtomicLevel

	*zlogger.Logger
	closers []io.Closer
	cores   map[string]zapcore.Core
	loggers map[string]*zlogger.Logger
}

func New(cfg Config, opts ...Option) (*Logger, error) {
	var logger Logger
	if err := logger.init(cfg, opts); err != nil {
		return nil, err
	}
	return &logger, nil
}

func (p *Logger) init(cfg Config, opts []Option) (err error) {
	p.level = zap.NewAtomicLevelAt(zbase.ZapLevel(cfg.Level))
	for _, apply := range opts {
		apply(p)
	}

	p.closers = make([]io.Closer, 0, len(cfg.Cores))
	p.cores = make(map[string]zapcore.Core)
	for _, core := range cfg.Cores {
		if err = p.openCore(core); err != nil {
			p.closeSinks()
			return fmt.Errorf("open core %q: %w", core.Name, err)
		}
	}

	p.loggers = make(map[string]*zlogger.Logger)
	for _, logger := range cfg.Loggers {
		if err = p.openLogger(logger); err != nil {
			p.closeSinks()
			return fmt.Errorf("open logger %q: %w", logger.Name, err)
		}
	}

	if len(cfg.Loggers) <= 0 {
		p.closeSinks()
		return errors.New("can't find any loggers")
	}

	name := cfg.Loggers[0].Name
	p.Logger = p.loggers[name]

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

	sink, err := newSinks(cfg.URLs)
	if err != nil {
		return fmt.Errorf("new sinks: %w", err)
	}

	enab := &levelEnabler{
		min:   zbase.ZapLevel(cfg.MinLevel),
		max:   zbase.ZapLevel(cfg.MaxLevel),
		level: p.level,
	}

	p.closers = append(p.closers, sink)
	p.cores[cfg.Name] = zapcore.NewCore(enc, sink, enab)

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

	opts := buildLoggerOptions(cfg)
	p.loggers[cfg.Name] = zlogger.New(cfg.Name, core, p.hook, opts...)

	return nil
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

func buildLoggerOptions(cfg LoggerConfig) []zap.Option {
	var opts []zap.Option

	if !cfg.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}

	switch cfg.StacktraceLevel {
	case PanicStacktrace:
		opts = append(opts, zap.AddStacktrace(zapcore.DPanicLevel))
	case ErrorStacktrace:
		opts = append(opts, zap.AddStacktrace(zapcore.ErrorLevel))
	case WarnStacktrace:
		opts = append(opts, zap.AddStacktrace(zapcore.WarnLevel))
	}

	return opts
}

func (p *Logger) closeSinks() {
	for _, c := range p.closers {
		c.Close()
	}
}

func (p *Logger) Close() (err error) {
	err = multierr.Append(err, p.Sync())
	for _, c := range p.closers {
		err = multierr.Append(err, c.Close())
	}
	return err
}

func (p *Logger) Sync() (err error) {
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

func (p *Logger) Named(name string) iface.Logger {
	if logger, ok := p.loggers[name]; ok {
		return logger
	}
	return p.Logger.Named(name)
}
