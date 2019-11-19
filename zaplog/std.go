package zaplog

import (
	"context"

	"github.com/ironzhang/tlog/iface"
	"go.uber.org/zap/zapcore"
)

var stdConfig = Config{
	Level: iface.DEBUG,
	Sinks: []SinkConfig{
		{
			Name: "StderrSink",
			URL:  "stderr",
		},
	},
	Cores: []CoreConfig{
		{
			Name:     "StderrCore",
			Encoding: "console",
			Encoder: EncoderConfig{
				TimeKey:        "T",
				LevelKey:       "L",
				NameKey:        "N",
				CallerKey:      "C",
				MessageKey:     "M",
				StacktraceKey:  "S",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    CapitalLevelEncoder,
				EncodeTime:     ISO8601TimeEncoder,
				EncodeDuration: StringDurationEncoder,
				EncodeCaller:   ShortCallerEncoder,
			},
			MinLevel: iface.DEBUG,
			MaxLevel: iface.FATAL,
			Sinks:    []string{"StderrSink"},
		},
	},
	Loggers: []LoggerConfig{
		{
			Name:            "default",
			DisableCaller:   false,
			StacktraceLevel: DisableStacktrace,
			Cores:           []string{"StderrCore"},
		},
	},
}

var stdLogger *Logger

func init() {
	hook := func(ctx context.Context) (args []interface{}) {
		return StdContextHook(ctx)
	}
	logger, err := New(stdConfig, ContextHookFunc(hook))
	if err != nil {
		panic(err)
	}
	stdLogger = logger
}

var StdContextHook = func(ctx context.Context) (args []interface{}) {
	return nil
}

func StdLogger() *Logger {
	return stdLogger
}
