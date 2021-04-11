package iface

import (
	"bytes"
	"errors"
	"fmt"
)

var errUnmarshalNilLevel = errors.New("can't unmarshal a nil *Level")

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

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "debug"
	case INFO:
		return "info"
	case WARN:
		return "warn"
	case ERROR:
		return "error"
	case PANIC:
		return "panic"
	case FATAL:
		return "fatal"
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}

func (l Level) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

func (l *Level) UnmarshalText(text []byte) error {
	if l == nil {
		return errUnmarshalNilLevel
	}
	if !l.unmarshalText(text) && !l.unmarshalText(bytes.ToLower(text)) {
		return fmt.Errorf("unrecognized level %q", text)
	}
	return nil
}

func (l *Level) unmarshalText(text []byte) bool {
	switch string(text) {
	case "debug", "DEBUG":
		*l = DEBUG
	case "info", "INFO", "":
		*l = INFO
	case "warn", "WARN":
		*l = WARN
	case "error", "ERROR":
		*l = ERROR
	case "panic", "PANIC":
		*l = PANIC
	case "fatal", "FATAL":
		*l = FATAL
	default:
		return false
	}
	return true
}

func StringToLevel(s string) (Level, error) {
	var l Level
	err := l.UnmarshalText([]byte(s))
	return l, err
}
