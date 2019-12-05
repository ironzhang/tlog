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

	iface.Logger
	loggers map[string]*zlogger.Logger
	closers []io.Closer
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

	p.loggers = make(map[string]*zlogger.Logger)
	for _, logger := range cfg.Loggers {
		if err = p.openLogger(logger); err != nil {
			p.closeLoggers()
			return err
		}
	}

	if len(cfg.Loggers) <= 0 {
		return errors.New("can't find any loggers")
	}
	name := cfg.Loggers[0].Name
	p.Logger = p.loggers[name]

	return nil
}

func (p *Logger) openLogger(cfg LoggerConfig) error {
	if _, ok := p.loggers[cfg.Name]; ok {
		return fmt.Errorf("logger %q is already opened", cfg.Name)
	}

	logger, closers, err := newLogger(p.level, p.hook, cfg)
	if err != nil {
		return fmt.Errorf("new logger %q: %w", cfg.Name, err)
	}

	p.loggers[cfg.Name] = logger
	p.closers = append(p.closers, closers...)
	return nil
}

func newLogger(level zap.AtomicLevel, hook ContextHook, cfg LoggerConfig) (*zlogger.Logger, []io.Closer, error) {
	opts := buildLoggerOptions(cfg)
	core, closers, err := buildLoggerCore(level, cfg)
	if err != nil {
		return nil, nil, err
	}
	return zlogger.New(cfg.Name, core, hook, opts...), closers, nil
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

func buildLoggerCore(level zap.AtomicLevel, cfg LoggerConfig) (zapcore.Core, []io.Closer, error) {
	enc, err := newEncoder(cfg.Encoding, cfg.Encoder)
	if err != nil {
		return nil, nil, fmt.Errorf("new encoder: %w", err)
	}

	cores := make([]zapcore.Core, 0, len(cfg.Outputs))
	closers := make([]io.Closer, 0, len(cfg.Outputs))
	closef := func() {
		for _, c := range closers {
			c.Close()
		}
	}

	for _, out := range cfg.Outputs {
		sink, err := newSinks(out.URLs)
		if err != nil {
			closef()
			return nil, nil, fmt.Errorf("new sinks: %w", err)
		}
		closers = append(closers, sink)

		enab := &levelEnabler{
			min:   zbase.ZapLevel(out.MinLevel),
			max:   zbase.ZapLevel(out.MaxLevel),
			level: level,
		}
		cores = append(cores, zapcore.NewCore(enc.Clone(), sink, enab))
	}
	return zapcore.NewTee(cores...), closers, nil
}

func (p *Logger) closeLoggers() {
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
	for _, logger := range p.loggers {
		err = multierr.Append(err, logger.Sync())
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
