package zconfig

import (
	"github.com/ironzhang/tlog/logger"
)

type Config struct {
	Level   logger.Level `json:"level" yaml:"level"`
	Sinks   []SinkConfig
	Cores   []CoreConfig   `json:"devices" yaml:"devices"`
	Loggers []LoggerConfig `json:"loggers" yaml:"loggers"`
}
