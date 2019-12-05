package zaplog

import (
	"context"

	"github.com/ironzhang/tlog/iface"
)

var stdConfig = Config{
	Level: iface.DEBUG,
	Loggers: []LoggerConfig{
		{
			Name:            "",
			DisableCaller:   false,
			StacktraceLevel: DisableStacktrace,
			Encoding:        "console",
			Encoder: EncoderConfig{
				TimeKey:        "T",
				LevelKey:       "L",
				NameKey:        "N",
				CallerKey:      "C",
				MessageKey:     "M",
				StacktraceKey:  "S",
				EncodeLevel:    CapitalLevelEncoder,
				EncodeTime:     ISO8601TimeEncoder,
				EncodeDuration: StringDurationEncoder,
				EncodeCaller:   ShortCallerEncoder,
			},
			Outputs: []OutputConfig{
				{
					MinLevel: iface.DEBUG,
					MaxLevel: iface.FATAL,
					URLs:     []string{"stderr"},
				},
			},
		},
	},
}

var stdLogger *Logger

func init() {
	hook := func(ctx context.Context) (args []interface{}) {
		return StdContextHook(ctx)
	}
	logger, err := New(stdConfig, SetContextHook(ContextHookFunc(hook)))
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
