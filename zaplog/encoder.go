package zaplog

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

func newEncoder(name string, enc zapcore.EncoderConfig) (zapcore.Encoder, error) {
	switch name {
	case "console":
		return zapcore.NewConsoleEncoder(enc), nil
	case "json":
		return zapcore.NewJSONEncoder(enc), nil
	}
	return nil, fmt.Errorf("no encoder for name %q", name)
}
