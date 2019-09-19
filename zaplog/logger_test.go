package zaplog

/*
func NewZapLogger(t *testing.T) *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	cfg.DisableStacktrace = true
	cfg.Level.SetLevel(TraceLevel)
	//l, err := cfg.Build(zap.AddStacktrace(zap.NewAtomicLevelAt(zap.DPanicLevel)))
	l, err := cfg.Build()
	if err != nil {
		t.Fatalf("build zap logger: %v", err)
	}
	return l
}

func TestLoggerSweetenFields(t *testing.T) {
	base := NewZapLogger(t)
	log := NewLogger(base)
	tests := []struct {
		args   []interface{}
		fields []zap.Field
	}{
		{
			args: []interface{}{"k1", 1, "k2", "2"},
			fields: []zap.Field{
				zap.Any("k1", 1),
				zap.Any("k2", "2"),
			},
		},
	}
	for i, tt := range tests {
		fields := log.sweetenFields(tt.args)
		if got, want := fields, tt.fields; !reflect.DeepEqual(got, want) {
			t.Errorf("%d: sweetenFields: got %v, want %v", i, got, want)
		} else {
			t.Logf("%d: sweetenFields: got %v", i, got)
		}
	}
}

func TestLoggerLog(t *testing.T) {
	base := NewZapLogger(t)
	lg := NewLogger(base)

	type Func func(...interface{})
	funcs := []Func{
		lg.Trace,
		lg.Debug,
		lg.Info,
		lg.Warn,
		lg.Error,
		//lg.Panic,
		//lg.Fatal,
	}
	for _, log := range funcs {
		log("hello, world")
	}
}

func TestLoggerLogf(t *testing.T) {
	base := NewZapLogger(t)
	lg := NewLogger(base)

	type Func func(string, ...interface{})
	funcs := []Func{
		lg.Tracef,
		lg.Debugf,
		lg.Infof,
		lg.Warnf,
		lg.Errorf,
		//lg.Panic,
		//lg.Fatal,
	}
	for _, log := range funcs {
		log("hello, world, function=%s", "TestLoggerLogf")
	}
}

func TestLoggerLogw(t *testing.T) {
	base := NewZapLogger(t)
	lg := NewLogger(base)

	type Func func(string, ...interface{})
	funcs := []Func{
		lg.Tracew,
		lg.Debugw,
		lg.Infow,
		lg.Warnw,
		lg.Errorw,
		//lg.Panic,
		//lg.Fatal,
	}
	for _, logw := range funcs {
		logw("hello, world", "function", "TestLoggerLogw")
	}
}

func TestLoggerWithArgs(t *testing.T) {
	base := NewZapLogger(t)
	log := NewLogger(base)
	log.WithArgs("function", "TestLoggerWithArgs").Debug("hello, world")
}

func TestLoggerWithContext(t *testing.T) {
	hook := func(context.Context) []interface{} {
		return []interface{}{"function", "TestLoggerWithContext"}
	}

	base := NewZapLogger(t)
	log := NewLogger(base, SetContextHook(hook))
	log.WithContext(context.Background()).Debug("hello, world")
	log.WithArgs("args", 1).WithContext(context.Background()).Debug("hello, world")
}

//func TestTee(t *testing.T) {
//	encoderCfg := zapcore.EncoderConfig{
//		MessageKey:     "msg",
//		LevelKey:       "level",
//		NameKey:        "logger",
//		EncodeLevel:    zapcore.LowercaseLevelEncoder,
//		EncodeTime:     zapcore.ISO8601TimeEncoder,
//		EncodeDuration: zapcore.StringDurationEncoder,
//	}
//
//	f1, _ := os.Create("./debug.log")
//	f2, _ := os.Create("./info.log")
//	core1 := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), f1, zap.DebugLevel)
//	core2 := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), f2, zap.InfoLevel)
//	core := zapcore.NewTee(core1, core2)
//	logger := zap.New(core)
//	logger.Debug("debug")
//	logger.Info("info")
//}
*/
