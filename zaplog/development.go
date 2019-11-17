package zaplog

import (
	"context"

	"github.com/ironzhang/tlog/iface"
	"go.uber.org/zap"
)

var DevelopmentConfig = Config{
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
			Name:  "default",
			Cores: []string{"StderrCore"},
		},
	},
}

var DevelopmentContextHook = func(ctx context.Context) (args []interface{}) {
	return nil
}

var DevelopmentLogger *Logger

func init() {
	hook := func(ctx context.Context) (args []interface{}) {
		return DevelopmentContextHook(ctx)
	}
	logger, err := New(DevelopmentConfig, ContextHookFunc(hook))
	if err != nil {
		panic(err)
	}
	DevelopmentLogger = logger
}
