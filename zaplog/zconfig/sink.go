package zconfig

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type SinkConfig struct {
	URL    string
	Params map[string]interface{}
}

func NewSink(cfg SinkConfig) (zap.Sink, error) {
	ws, cf, err := zap.Open(cfg.URL)
	if err != nil {
		return nil, err
	}
	return &sink{
		WriteSyncer: ws,
		closef:      cf,
	}, nil
}

type sink struct {
	zapcore.WriteSyncer
	closef func()
}

func (s *sink) Close() error {
	s.closef()
	return nil
}
