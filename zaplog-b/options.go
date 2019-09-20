package zaplog

import "github.com/ironzhang/tlog"

type Option func(*Logger)

func SetContextHook(h tlog.ContextHookFunc) Option {
	return func(l *Logger) {
		l.hook = h
	}
}
