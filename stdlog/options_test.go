package stdlog

import (
	"log"
	"os"
	"testing"
)

func TestSetLevelOption(t *testing.T) {
	base := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	tests := []Level{
		DEBUG,
		TRACE,
		INFO,
		WARN,
		ERROR,
		PANIC,
		FATAL,
		-10,
		100,
	}
	for i, lv := range tests {
		lg := NewLogger(base, SetLevel(lv))
		if got, want := lg.GetLevel(), lv; got != want {
			t.Errorf("%d: level: got %v, want %v", i, got, want)
		} else {
			t.Logf("%d: level: got %v", i, got)
		}
	}
}

func TestSetCalldepthOption(t *testing.T) {
	base := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	for i := 0; i < 10; i++ {
		lg := NewLogger(base, SetCalldepth(i))
		if got, want := lg.GetCalldepth(), i; got != want {
			t.Errorf("%d: call depth: got %v, want %v", i, got, want)
		} else {
			t.Logf("%d: call depth: got %v", i, got)
		}
	}
}
