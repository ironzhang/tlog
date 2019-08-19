package stdlog

import "github.com/ironzhang/tlog/logger"

type Option func(*Logger)

func SetLevel(lv Level) Option {
	return func(l *Logger) {
		l.SetLevel(lv)
	}
}

func SetCalldepth(calldepth int) Option {
	return func(l *Logger) {
		l.SetCalldepth(calldepth)
	}
}

func SetContextHook(hook logger.ContextHookFunc) Option {
	return func(l *Logger) {
		l.SetContextHook(hook)
	}
}
