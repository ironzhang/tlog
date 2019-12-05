package zaplog

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

func newEncoder(name string, cfg EncoderConfig) (zapcore.Encoder, error) {
	enc := zapcore.EncoderConfig{
		MessageKey:     cfg.MessageKey,
		LevelKey:       cfg.LevelKey,
		TimeKey:        cfg.TimeKey,
		NameKey:        cfg.NameKey,
		CallerKey:      cfg.CallerKey,
		StacktraceKey:  cfg.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    cfg.EncodeLevel.zap(),
		EncodeTime:     cfg.EncodeTime.zap(),
		EncodeDuration: cfg.EncodeDuration.zap(),
		EncodeCaller:   cfg.EncodeCaller.zap(),
		EncodeName:     cfg.EncodeName.zap(),
	}
	if enc.MessageKey == "" {
		enc.MessageKey = "M"
	}
	if enc.StacktraceKey == "" {
		enc.StacktraceKey = "S"
	}

	switch name {
	case "console", "":
		return zapcore.NewConsoleEncoder(enc), nil
	case "json":
		return zapcore.NewJSONEncoder(enc), nil
	}
	return nil, fmt.Errorf("no encoder for name %q", name)
}
