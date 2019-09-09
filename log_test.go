package tlog_test

import (
	"github.com/ironzhang/tlog"
)

//func TestLog(t *testing.T) {
//	//hook := func(ctx context.Context) []interface{} {
//	//	return []interface{}{"TraceID", "123456"}
//	//}
//	//l := zaplog.Std.WithOptions(zaplog.SetContextHook(hook))
//	//tlog.SetLogger(l)
//
//	tlog.Trace("hello")
//	tlog.Debug("hello")
//	tlog.Info("hello")
//	tlog.WithArgs("function", "TestLog").Warn("hello")
//	tlog.WithContext(context.Background()).WithArgs("function", "TestLog").Error("hello")
//}

func Example_log() {
	tlog.Debug("debug")

	// Output:
	// debug
}
