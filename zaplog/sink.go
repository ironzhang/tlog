package zaplog

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newSink(name, url string) (*sink, error) {
	ws, cf, err := zap.Open(url)
	if err != nil {
		return nil, fmt.Errorf("zap open sink %q: %w", name, err)
	}
	return &sink{
		name:        name,
		closef:      cf,
		WriteSyncer: ws,
	}, nil
}

type sink struct {
	name   string
	closef func()
	zapcore.WriteSyncer
}

func (s *sink) Close() error {
	s.closef()
	return nil
}
