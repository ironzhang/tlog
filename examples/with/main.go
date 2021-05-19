package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ironzhang/tlog"
	"github.com/ironzhang/tlog/zaplog"
)

func ContextHook(ctx context.Context) (args []interface{}) {
	traceID, ok := ctx.Value("trace_id").(string)
	if ok {
		args = append(args, "trace_id", traceID)
	}
	return args
}

func main() {
	cfg := zaplog.NewDevelopmentConfig()
	logger, err := zaplog.New(cfg, zaplog.SetContextHook(zaplog.ContextHookFunc(ContextHook)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "new logger: %v\n", err)
		return
	}
	defer logger.Close()

	tlog.SetLogger(logger)

	ctx := context.WithValue(context.Background(), "trace_id", "123456")
	PrintLogWithContext(ctx)
	PrintLogWithArgs(ctx)
}

func PrintLogWithContext(ctx context.Context) {
	tlog.WithContext(ctx).Info("print log with context")
}

func PrintLogWithArgs(ctx context.Context) {
	log := tlog.WithArgs("a1", 1, "a2", 2)
	log.Info("print log with args")
}
