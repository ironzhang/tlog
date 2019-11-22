package zaplog

type Option func(*Logger)

func SetContextHook(h ContextHook) Option {
	return func(p *Logger) {
		p.hook = h
	}
}
