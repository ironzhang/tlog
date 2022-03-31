package main

import (
	"context"
	"fmt"
	"os"

	"git.xiaojukeji.com/pearls/tlog"
	"git.xiaojukeji.com/pearls/tlog/zaplog"
)

func ContextHook(ctx context.Context) (args []interface{}) {
	traceID, ok := ctx.Value("trace_id").(string)
	if ok {
		args = append(args, "trace_id", traceID)
	}
	return args
}

func main() {
	cfg := zaplog.NewDevelopmentConfig() // 构造日志配置
	//cfg := zaplog.NewProductionConfig()                                                        // 构造日志配置
	logger, err := zaplog.New(cfg, zaplog.SetContextHook(zaplog.ContextHookFunc(ContextHook))) // 构造日志对象
	if err != nil {
		fmt.Fprintf(os.Stderr, "new logger: %v\n", err)
		return
	}
	defer logger.Close()   // 退出程序前关闭日志对象
	tlog.SetLogger(logger) // 设置日志对象

	ctx := context.WithValue(context.Background(), "trace_id", "123456")
	tlog.WithContext(ctx).Info("hello, world") // 输出日志
}
