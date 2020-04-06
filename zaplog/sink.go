package zaplog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	_ "git.xiaojukeji.com/pearls/tlog/zaplog/zsink"
)

func newSinks(urls []string) (zap.Sink, error) {
	ws, cf, err := zap.Open(urls...)
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
