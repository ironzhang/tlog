package stdlog

import (
	"log"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	log := NewLogger(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile), DEBUG, 0)
	log.Debug("hello, world")
	log.Debugf("hello, %s", "china")
	log.Debugw("hello beijing", "LogLevel", DEBUG)
}
