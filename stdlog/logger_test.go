package stdlog

import (
	"context"
	"io"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/ironzhang/tlog/logger"
)

func NewBaseLogger() *log.Logger {
	var w io.Writer
	w = os.Stderr
	//w = ioutil.Discard
	return log.New(w, "", log.LstdFlags|log.Lshortfile)
}

func TestLoggerOutput(t *testing.T) {
	base := NewBaseLogger()
	//base.SetOutput(ioutil.Discard)
	log := NewLogger(base)
	log.SetLevel(TRACE)

	for lv := DEBUG; lv <= FATAL; lv++ {
		log.Output(lv, 1, "hello", "world")
	}
	for lv := DEBUG; lv <= FATAL; lv++ {
		log.Outputf(lv, 1, "%s hello, world", lv)
	}
	for lv := DEBUG; lv <= FATAL; lv++ {
		log.Outputw(lv, 1, "hello, world", "level", lv)
	}
}

func TestLogger(t *testing.T) {
	base := NewBaseLogger()
	//base.SetOutput(ioutil.Discard)
	lg := NewLogger(base)
	lg.SetLevel(DEBUG)

	type LogFunc func(...interface{})
	logFuncs := []LogFunc{
		lg.Debug,
		lg.Trace,
		lg.Info,
		lg.Warn,
		lg.Error,
		lg.Panic,
		lg.Fatal,
	}
	for _, log := range logFuncs {
		log("hello, world")
	}

	type LogfFunc func(string, ...interface{})
	logfFuncs := []LogfFunc{
		lg.Debugf,
		lg.Tracef,
		lg.Infof,
		lg.Warnf,
		lg.Errorf,
		lg.Panicf,
		lg.Fatalf,
	}
	for _, logf := range logfFuncs {
		logf("hello, world. function=%s", "TestLogger")
	}

	type LogwFunc func(string, ...interface{})
	logwFuncs := []LogwFunc{
		lg.Debugw,
		lg.Tracew,
		lg.Infow,
		lg.Warnw,
		lg.Errorw,
		lg.Panicw,
		lg.Fatalw,
	}
	for _, logw := range logwFuncs {
		logw("hello, world", "function", "TestLogger")
		logw("hello, world", 100, "TestLogger")
		logw("hello, world", 100, "TestLogger", "last")
	}
}

func TestLoggerWithArgs(t *testing.T) {
	base := NewBaseLogger()
	//base.SetOutput(ioutil.Discard)
	lg := NewLogger(base, SetCalldepth(1))
	lg.SetLevel(DEBUG)
	lg.Debugw("hello, TestLoggerWithArgs")
	lg.WithArgs("function", "TestLoggerWithArgs").Debugw("hello, TestLoggerWithArgs")
}

func TestLoggerWithContext(t *testing.T) {
	logger.WithContextHook = func(ctx context.Context) []interface{} {
		return []interface{}{"TraceID", "123466", "SpanID", "1"}
	}

	base := NewBaseLogger()
	//base.SetOutput(ioutil.Discard)
	lg := NewLogger(base, SetCalldepth(1))
	lg.SetLevel(DEBUG)
	lg.Debugw("hello, TestLoggerWithContext")
	lg.WithContext(context.Background()).Debugw("hello, TestLoggerWithContext")
}

func TestSweetenFields(t *testing.T) {
	tests := []struct {
		args   []interface{}
		fields []field
	}{
		{
			args: []interface{}{"n", 1},
			fields: []field{
				{key: "n", value: 1},
			},
		},
	}
	for i, tt := range tests {
		fields := sweetenFields(tt.args)
		if got, want := fields, tt.fields; !reflect.DeepEqual(got, want) {
			t.Errorf("%d: fields: got %v, want %v", i, got, want)
		} else {
			t.Logf("%d: fields: got %v", i, got)
		}
	}
}
