package tlog_test

import (
	"testing"

	"github.com/ironzhang/tlog"
)

func TestLog(t *testing.T) {
	//	w := &lumberjack.Logger{
	//		Filename: "./test.log",
	//	}
	//	defer w.Close()
	//	tlog.Default.SetLogger(log.New(w, "", log.LstdFlags|log.Llongfile))

	tlog.Debug("hello")
	tlog.Info("hello")
	tlog.WithArgs("function", "TestLog").Warn("hello")
}
