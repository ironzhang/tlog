package zaplog

import (
	"go.uber.org/zap"
)

var Std *Logger

func init() {
	cfg := zap.NewDevelopmentConfig()
	cfg.DisableStacktrace = true
	cfg.Level.SetLevel(zap.DebugLevel)
	base, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	Std = NewLogger(base)
}
