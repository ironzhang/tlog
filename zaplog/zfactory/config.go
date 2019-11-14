package zfactory

import (
	"github.com/ironzhang/tlog/logger"
	"go.uber.org/zap/zapcore"
)

type Level = logger.Level

type Encoder struct {
	Encoding       string                  `json:"encoding" yaml:"encoding"`
	MessageKey     string                  `json:"messageKey" yaml:"messageKey"`
	LevelKey       string                  `json:"levelKey" yaml:"levelKey"`
	TimeKey        string                  `json:"timeKey" yaml:"timeKey"`
	NameKey        string                  `json:"nameKey" yaml:"nameKey"`
	CallerKey      string                  `json:"callerKey" yaml:"callerKey"`
	StacktraceKey  string                  `json:"stacktraceKey" yaml:"stacktraceKey"`
	LineEnding     string                  `json:"lineEnding" yaml:"lineEnding"`
	EncodeLevel    zapcore.LevelEncoder    `json:"levelEncoder" yaml:"levelEncoder"`
	EncodeTime     zapcore.TimeEncoder     `json:"timeEncoder" yaml:"timeEncoder"`
	EncodeDuration zapcore.DurationEncoder `json:"durationEncoder" yaml:"durationEncoder"`
	EncodeCaller   zapcore.CallerEncoder   `json:"callerEncoder" yaml:"callerEncoder"`
	EncodeName     zapcore.NameEncoder     `json:"nameEncoder" yaml:"nameEncoder"`
}

type Device struct {
	Name     string   `json:"name" yaml:"name"`
	Encoder  Encoder  `json:"encoder" yaml:"encoder"`
	MinLevel Level    `json:"minLevel" yaml:"minLevel"`
	MaxLevel Level    `json:"maxLevel" yaml:"maxLevel"`
	URLs     []string `json:"urls" yaml:"urls"`
}

type Logger struct {
	Name    string   `json:"name" yaml:"name"`
	Level   Level    `json:"level" yaml:"level"`
	Devices []string `json:"devices" yaml:"devices"`
}

type Config struct {
	Devices []Device `json:"devices" yaml:"devices"`
	Loggers []Logger `json:"loggers" yaml:"loggers"`
}
