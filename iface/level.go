package iface

// 日志级别
type Level int8

// 日志级别常量定义
const (
	DEBUG Level = iota - 1
	INFO
	WARN
	ERROR
	PANIC
	FATAL
)
