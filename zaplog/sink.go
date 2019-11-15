package zaplog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type closedSink struct {
	zapcore.WriteSyncer
	closef func() error
}

func (p *closedSink) Close() error {
	return p.closef()
}

func wrapClosedSink(syncer zapcore.WriteSyncer, closef func() error) zap.Sink {
	return &closedSink{
		WriteSyncer: syncer,
		closef:      closef,
	}
}
