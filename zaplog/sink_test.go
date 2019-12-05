package zaplog

import (
	"net/url"
	"testing"

	"go.uber.org/zap"
)

type tSink struct {
	writeCount int
	syncCount  int
	closeCount int
}

func (p *tSink) Write(b []byte) (int, error) {
	p.writeCount++
	return 0, nil
}

func (p *tSink) Sync() error {
	p.syncCount++
	return nil
}

func (p *tSink) Close() error {
	p.closeCount++
	return nil
}

func RegisterTestSink(t testing.TB, scheme string) *tSink {
	sink := tSink{}
	err := zap.RegisterSink(scheme, func(u *url.URL) (zap.Sink, error) {
		return &sink, nil
	})
	if err != nil {
		t.Fatalf("register sink: %v", err)
	}
	return &sink
}

func TestSink(t *testing.T) {
	tsink := RegisterTestSink(t, "TestSink")

	urls := []string{"TestSink://1", "TestSink://2"}
	sink, err := newSinks(urls)
	if err != nil {
		t.Fatalf("new sinks: %v", err)
	}
	sink.Write([]byte{})
	sink.Sync()
	sink.Close()

	if got, want := tsink.writeCount, len(urls); got != want {
		t.Errorf("write: got %v, want %v", got, want)
	}
	if got, want := tsink.syncCount, len(urls); got != want {
		t.Errorf("sync: got %v, want %v", got, want)
	}
	if got, want := tsink.closeCount, len(urls); got != want {
		t.Errorf("close: got %v, want %v", got, want)
	}
	t.Logf("writeCount: %d, syncCount: %d, closeCount: %d", tsink.writeCount, tsink.syncCount, tsink.closeCount)
}
