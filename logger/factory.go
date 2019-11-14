package logger

type Factory interface {
	GetLogger(name string) Logger
}
