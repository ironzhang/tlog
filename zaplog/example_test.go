package zaplog

func ExampleStdLogger() {
	logger := StdLogger()
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")
	//logger.Panic("panic")

	// output:
}