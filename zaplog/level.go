package zaplog

type Level int8

const (
	TRACE Level = iota - 2
	DEBUG
	INFO
	WARN
	ERROR
	PANIC
	FATAL
)
