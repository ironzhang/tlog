package zaplog

import (
	"go.uber.org/zap/zapcore"

	"github.com/ironzhang/tlog/iface"
)

type SinkConfig struct {
	Name string
	URL  string
}

type CoreConfig struct {
	Name     string                `json:"name" yaml:"name"`
	Encoding string                `json:"encoding" yaml:"encoding"`
	Encoder  zapcore.EncoderConfig `json:"encoder" yaml:"encoder"`
	MinLevel iface.Level           `json:"minLevel" yaml:"minLevel"`
	MaxLevel iface.Level           `json:"maxLevel" yaml:"maxLevel"`
	Sinks    []string              `json:"urls" yaml:"urls"`
}

type LoggerConfig struct {
	Name  string   `json:"name" yaml:"name"`
	Cores []string `json:"devices" yaml:"devices"`
}

type Config struct {
	Level   iface.Level `json:"level" yaml:"level"`
	Sinks   []SinkConfig
	Cores   []CoreConfig   `json:"devices" yaml:"devices"`
	Loggers []LoggerConfig `json:"loggers" yaml:"loggers"`
}
