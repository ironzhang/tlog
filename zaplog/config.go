package zaplog

import (
	"github.com/ironzhang/tlog/logger"
	"go.uber.org/zap/zapcore"
)

type SinkConfig struct {
	Name string
	URL  string
}

type CoreConfig struct {
	Name     string                `json:"name" yaml:"name"`
	Encoding string                `json:"encoding" yaml:"encoding"`
	Encoder  zapcore.EncoderConfig `json:"encoder" yaml:"encoder"`
	MinLevel logger.Level          `json:"minLevel" yaml:"minLevel"`
	MaxLevel logger.Level          `json:"maxLevel" yaml:"maxLevel"`
	SinkRefs []string              `json:"urls" yaml:"urls"`
}

type LoggerConfig struct {
	Name     string   `json:"name" yaml:"name"`
	CoreRefs []string `json:"devices" yaml:"devices"`
}

type Config struct {
	Level   logger.Level   `json:"level" yaml:"level"`
	Cores   []CoreConfig   `json:"devices" yaml:"devices"`
	Loggers []LoggerConfig `json:"loggers" yaml:"loggers"`
}
