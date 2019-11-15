package logger

type Factory interface {
	GetDefaultLogger() Logger
	GetLogger(name string) Logger
}
