package zaplog

import (
	"go.uber.org/zap/zapcore"

	"github.com/ironzhang/tlog/iface"
)

type SinkConfig struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url" yaml:"url"`
}

type CoreConfig struct {
	Name     string                `json:"name" yaml:"name"`
	Encoding string                `json:"encoding" yaml:"encoding"`
	Encoder  zapcore.EncoderConfig `json:"encoder" yaml:"encoder"`
	MinLevel iface.Level           `json:"minLevel" yaml:"minLevel"`
	MaxLevel iface.Level           `json:"maxLevel" yaml:"maxLevel"`
	Sinks    []string              `json:"sinks" yaml:"sinks"`
}

type LoggerConfig struct {
	Name  string   `json:"name" yaml:"name"`
	Cores []string `json:"cores" yaml:"cores"`
}

type Config struct {
	Level   iface.Level    `json:"level" yaml:"level"`
	Sinks   []SinkConfig   `json:"sinks" yaml:"sinks"`
	Cores   []CoreConfig   `json:"cores" yaml:"cores"`
	Loggers []LoggerConfig `json:"loggers" yaml:"loggers"`
}
