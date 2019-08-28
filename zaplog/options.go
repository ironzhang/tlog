package zaplog

import "github.com/ironzhang/tlog/logger"

type Option func(*Logger)

func SetContextHook(h logger.ContextHookFunc) Option {
	return func(l *Logger) {
		l.hook = h
	}
}
