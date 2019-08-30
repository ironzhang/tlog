package zapx

import (
	"fmt"
	"net/url"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Sink interface {
	Name() string
	zap.Sink
}

type sinkList []Sink

type sink struct {
	name string
	zap.Sink
}

func newSink(name string, rawURL string) (Sink, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("can not parse %q as a URL: %v", rawURL, err)
	}
	s, err := newZapSink(u)
	if err != nil {
		return nil, fmt.Errorf("can not new a zap sink: %v", err)
	}
	return &sink{name: name, Sink: s}, nil
}

func (s *sink) Name() string {
	return s.name
}

func newZapSink(u *url.URL) (zap.Sink, error) {
	switch u.Scheme {
	case "file":
		return newFileSink(u)
	default:
		return newFileSink(u)
	}
}

func newFileSink(u *url.URL) (zap.Sink, error) {
	if u.User != nil {
		return nil, fmt.Errorf("user and password not allowed with file URLs: got %v", u)
	}
	if u.Fragment != "" {
		return nil, fmt.Errorf("fragments not allowed with file URLs: got %v", u)
	}
	if u.RawQuery != "" {
		return nil, fmt.Errorf("query parameters not allowed with file URLs: got %v", u)
	}
	// Error messages are better if we check hostname and port separately.
	if u.Port() != "" {
		return nil, fmt.Errorf("ports not allowed with file URLs: got %v", u)
	}
	if hn := u.Hostname(); hn != "" && hn != "localhost" {
		return nil, fmt.Errorf("file URLs must leave host empty or use localhost: got %v", u)
	}
	switch u.Path {
	case "stdout":
		return nopCloserSink{os.Stdout}, nil
	case "stderr":
		return nopCloserSink{os.Stderr}, nil
	}
	return os.OpenFile(u.Path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
}

type nopCloserSink struct{ zapcore.WriteSyncer }

func (nopCloserSink) Close() error { return nil }
