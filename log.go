package tlog

var logging Logger = &nopLogger{}

var (
	WithArgs    = logging.WithArgs
	WithContext = logging.WithContext

	Trace  = logging.Trace
	Tracef = logging.Tracef
	Tracew = logging.Tracew

	Debug  = logging.Debug
	Debugf = logging.Debugf
	Debugw = logging.Debugw

	Info  = logging.Info
	Infof = logging.Infof
	Infow = logging.Infow

	Warn  = logging.Warn
	Warnf = logging.Warnf
	Warnw = logging.Warnw

	Error  = logging.Error
	Errorf = logging.Errorf
	Errorw = logging.Errorw

	Panic  = logging.Panic
	Panicf = logging.Panicf
	Panicw = logging.Panicw

	Fatal  = logging.Fatal
	Fatalf = logging.Fatalf
	Fatalw = logging.Fatalw
)

func SetLogger(l Logger) Logger {
	prev := logging
	if l == nil {
		logging = &nopLogger{}
	} else {
		logging = l
	}

	WithArgs = logging.WithArgs
	WithContext = logging.WithContext
	Trace = logging.Trace
	Tracef = logging.Tracef
	Tracew = logging.Tracew
	Debug = logging.Debug
	Debugf = logging.Debugf
	Debugw = logging.Debugw
	Info = logging.Info
	Infof = logging.Infof
	Infow = logging.Infow
	Warn = logging.Warn
	Warnf = logging.Warnf
	Warnw = logging.Warnw
	Error = logging.Error
	Errorf = logging.Errorf
	Errorw = logging.Errorw
	Panic = logging.Panic
	Panicf = logging.Panicf
	Panicw = logging.Panicw
	Fatal = logging.Fatal
	Fatalf = logging.Fatalf
	Fatalw = logging.Fatalw

	return prev
}
