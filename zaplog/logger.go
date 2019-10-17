package zaplog

import "github.com/ironzhang/tlog/logger"

type Level = logger.Level

const (
	DEBUG = logger.DEBUG
	INFO  = logger.INFO
	WARN  = logger.WARN
	ERROR = logger.ERROR
	PANIC = logger.PANIC
	FATAL = logger.FATAL
)

type Logger struct {
}

func NewLogger(cfg Config) (*Logger, error) {
	return &Logger{}, nil
}

func (p *Logger) init(cfg Config) (*Logger, error) {
	return p, nil
}
