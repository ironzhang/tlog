package stdlog

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestLevelString(t *testing.T) {
	tests := []struct {
		lv Level
		s  string
	}{
		{lv: DEBUG, s: "DEBUG"},
		{lv: TRACE, s: "TRACE"},
		{lv: INFO, s: "INFO"},
		{lv: WARN, s: "WARN"},
		{lv: ERROR, s: "ERROR"},
		{lv: PANIC, s: "PANIC"},
		{lv: FATAL, s: "FATAL"},
		{lv: -10, s: "UNKNOWN(-10)"},
		{lv: 100, s: "UNKNOWN(100)"},
	}
	for i, tt := range tests {
		if got, want := tt.lv.String(), tt.s; got != want {
			t.Errorf("%d: level: got %v, want %v", i, got, want)
		} else {
			t.Logf("%d: level: got %v", i, got)
		}
	}
}

func TestAtomicLevelLoad(t *testing.T) {
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
		atomic := newAtomicLevel(lv)
		if got, want := atomic.Load(), lv; got != want {
			t.Errorf("%d: atomic level: got %v, want %v", i, got, want)
		} else {
			t.Logf("%d: atomic level: got %v", i, got)
		}
	}
}

func TestAtomicLevelStore(t *testing.T) {
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
		atomic := newAtomicLevel(0)
		atomic.Store(lv)
		if got, want := Level(atomic.value), lv; got != want {
			t.Errorf("%d: atomic level: got %v, want %v", i, got, want)
		} else {
			t.Logf("%d: atomic level: got %v", i, got)
		}
	}
}

func TestAtomicLevelStoreAndLoad(t *testing.T) {
	seed := time.Now().Unix()
	//seed = 0
	rand.Seed(seed)

	var wg sync.WaitGroup
	atomic := newAtomicLevel(0)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			atomic.Store(Level(rand.Int31n(1000)))
			wg.Done()
		}()
	}
	wg.Wait()
	t.Logf("atomic level: load: %d", atomic.Load())

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
		atomic.Store(lv)
		if got, want := atomic.Load(), lv; got != want {
			t.Errorf("%d: atomic level: got %v, want %v", i, got, want)
		} else {
			t.Logf("%d: atomic level: got %v", i, got)
		}
	}
}
