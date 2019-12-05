package zaplog

import (
	"context"
)

var stdLogger *Logger

func init() {
	hook := func(ctx context.Context) (args []interface{}) {
		return StdContextHook(ctx)
	}
	logger, err := New(NewDevelopmentConfig(), SetContextHook(ContextHookFunc(hook)))
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
