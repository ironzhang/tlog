package zapx

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

func newEncoder(name string, encoderConfig zapcore.EncoderConfig) (zapcore.Encoder, error) {
	switch name {
	case "console":
		return zapcore.NewConsoleEncoder(encoderConfig), nil
	case "json":
		return zapcore.NewJSONEncoder(encoderConfig), nil
	default:
		return nil, fmt.Errorf("encoder %q is not supported", name)
	}
}
