package zaplog

import (
	"net/url"
	"testing"

	"go.uber.org/zap"
)

var (
	tSinkWrite int
	tSinkSync  int
	tSinkClose int
)

type tSink struct {
}

func (p *tSink) Write(b []byte) (int, error) {
	tSinkWrite++
	return 0, nil
}

func (p *tSink) Sync() error {
	tSinkSync++
	return nil
}

func (p *tSink) Close() error {
	tSinkClose++
	return nil
}

func RegisterTestSinks(t testing.TB) {
	err := zap.RegisterSink("tsink", func(u *url.URL) (zap.Sink, error) {
		return &tSink{}, nil
	})
	if err != nil {
		t.Fatalf("register sink: %v", err)
	}
}

func TestSink(t *testing.T) {
	RegisterTestSinks(t)

	urls := []string{"tsink://1", "tsink://2"}
	sink, err := newSinks(urls)
	if err != nil {
		t.Fatalf("new sinks: %v", err)
	}
	sink.Write([]byte{})
	sink.Sync()
	sink.Close()

	if got, want := tSinkWrite, len(urls); got != want {
		t.Errorf("write: got %v, want %v", got, want)
	}
	if got, want := tSinkSync, len(urls); got != want {
		t.Errorf("sync: got %v, want %v", got, want)
	}
	if got, want := tSinkClose, len(urls); got != want {
		t.Errorf("close: got %v, want %v", got, want)
	}
	t.Logf("write: %d, sync: %d, close: %d", tSinkWrite, tSinkSync, tSinkClose)
}
