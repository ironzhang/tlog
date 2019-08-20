package zaplog

import "go.uber.org/zap"

type Logger struct {
}

func NewLogger(logger *zap.Logger, lv zap.AtomicLevel) *Logger {
	return &Logger{}
}
