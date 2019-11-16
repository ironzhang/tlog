package zconfig

import (
	"github.com/ironzhang/tlog/logger"
	"go.uber.org/zap/zapcore"
)

type CoreConfig struct {
	Name     string                `json:"name" yaml:"name"`
	Encoding string                `json:"encoding" yaml:"encoding"`
	Encoder  zapcore.EncoderConfig `json:"encoder" yaml:"encoder"`
	MinLevel logger.Level          `json:"minLevel" yaml:"minLevel"`
	MaxLevel logger.Level          `json:"maxLevel" yaml:"maxLevel"`
	Sinks    []string              `json:"urls" yaml:"urls"`
}

func NewCore(cfg CoreConfig) (zapcore.Core, error) {
	return nil, nil
}
