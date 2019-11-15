package zaplog

import (
	"fmt"
	"net/url"
	"strings"

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

func (p *ZapLogger) buildCore(core zconfig.Core) error {
	return nil
}

func (p *ZapLogger) buildSink(rawurl string) error {
	key, err := parseKeyFromURL(rawurl)
	if err != nil {
		return err
	}
	if _, ok := p.sinks[key]; ok {
		return fmt.Errorf("sink %q is opened", key)
	}
	writer, closef, err := zap.Open(rawurl)
	if err != nil {
		return err
	}
	sink := wrapClosedSink(writer, func() error {
		closef()
		return nil
	})
	p.sinks[key] = sink
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

func parseKeyFromURL(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", fmt.Errorf("parse url: %w", err)
	}
	if u.Scheme == "" {
		u.Scheme = "file"
	}
	s := u.String()
	i := strings.IndexByte(s, '?')
	if i > 0 {
		s = s[:i]
	}
	return s, nil
}
