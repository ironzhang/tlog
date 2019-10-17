package logger

type Level int8

const (
	DEBUG Level = iota - 1
	INFO
	WARN
	ERROR
	PANIC
	FATAL
)
