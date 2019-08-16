package stdlog

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	var w io.Writer
	w = os.Stderr
	w = ioutil.Discard
	log := NewLogger(log.New(w, "", log.LstdFlags|log.Lshortfile), SetLevel(DEBUG))
	log.Debug("hello, world")
	log.Debugf("hello, %s", "china")
	log.Debugw("hello beijing", "LogLevel", DEBUG)
}
