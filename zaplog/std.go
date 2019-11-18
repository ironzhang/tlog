package zaplog

import (
	"context"

	"github.com/ironzhang/tlog/iface"
	"go.uber.org/zap"
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
			Encoder:  zap.NewDevelopmentEncoderConfig(),
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
