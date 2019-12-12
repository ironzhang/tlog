package zaplog

import (
	"fmt"
	"os"
)

func ExampleStdLogger() {
	logger := StdLogger()
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")
	//logger.Panic("panic")

	// output:
}

func ExampleNew() {
	logger, err := New(NewProductionConfig())
	if err != nil {
		fmt.Fprintf(os.Stderr, "new: %v", err)
		return
	}
	defer logger.Close()

	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")

	// output:
}
