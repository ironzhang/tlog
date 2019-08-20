package tlog_test

import (
	"log"
	"testing"

	"github.com/ironzhang/tlog"
	"github.com/natefinch/lumberjack"
)

func TestLog(t *testing.T) {
	w := &lumberjack.Logger{
		Filename: "./test.log",
	}
	defer w.Close()
	tlog.Default.SetLogger(log.New(w, "", log.LstdFlags|log.Llongfile))

	tlog.Debug("hello")
	tlog.Info("hello")
	tlog.WithArgs("function", "TestLog").Warn("hello")
}
