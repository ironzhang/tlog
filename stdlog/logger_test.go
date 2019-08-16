package stdlog

import (
	"io"
	"log"
	"os"
	"testing"
)

var base *log.Logger

func TestMain(t *testing.M) {
	var w io.Writer
	w = os.Stderr
	//w = ioutil.Discard
	base = log.New(w, "", log.LstdFlags|log.Lshortfile)
	t.Run()
}

func TestLoggerArgs(t *testing.T) {
	type PrintFunc func(args ...interface{})
	log := NewLogger(base)

	levels := []Level{
		DEBUG,
		TRACE,
		INFO,
		WARN,
		ERROR,
		PANIC,
		FATAL,
	}
	for _, lv := range levels {
		log.SetLevel(lv)
		Prints := []PrintFunc{
			log.Debug,
			log.Trace,
			log.Info,
			log.Warn,
			log.Error,
			log.Panic,
			log.Fatal,
		}
		for _, Print := range Prints {
			Print(lv, "\t", "hello, world")
			Print(lv, "\t", "hello, world", 0)
			Print(lv, "\t", "hello, world", " ", "hello, 世界", " ", 2)
		}
	}
}
