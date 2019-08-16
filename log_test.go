package tlog_test

import (
	"testing"

	"github.com/ironzhang/tlog"
)

func TestLog(t *testing.T) {
	tlog.Debug("hello")
	tlog.Info("hello")
	tlog.WithArgs("function", "TestLog").Warn("hello")
}
