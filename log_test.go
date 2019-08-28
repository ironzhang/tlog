package tlog_test

import (
	"testing"

	"github.com/ironzhang/tlog"
	"github.com/ironzhang/tlog/zaplog"
	"go.uber.org/zap"
)

func NewZapLogger(t *testing.T) *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	cfg.DisableStacktrace = true
	cfg.Level.SetLevel(zap.DebugLevel)
	//l, err := cfg.Build(zap.AddStacktrace(zap.NewAtomicLevelAt(zap.DPanicLevel)))
	l, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		t.Fatalf("build zap logger: %v", err)
	}
	return l
}

func TestLog(t *testing.T) {
	//	w := &lumberjack.Logger{
	//		Filename: "./test.log",
	//	}
	//	defer w.Close()
	//	tlog.Default.SetLogger(log.New(w, "", log.LstdFlags|log.Llongfile))

	base := NewZapLogger(t)
	log := zaplog.NewLogger(base)
	tlog.SetLogger(log)

	tlog.Debug("hello")
	tlog.Info("hello")
	tlog.WithArgs("function", "TestLog").Warn("hello")
}
