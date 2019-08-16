package tlog_test

import (
	"log"
	"os"
	"testing"

	"github.com/ironzhang/tlog"
)

func TestLog(t *testing.T) {
	tlog.Default.SetLogger(log.New(os.Stderr, "", log.LstdFlags|log.Llongfile))

	tlog.Debug("hello")
	tlog.Info("hello")
	tlog.WithArgs("function", "TestLog").Warn("hello")
}
