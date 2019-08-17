package stdlog

import (
	"io"
	"log"
	"os"
	"testing"
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
