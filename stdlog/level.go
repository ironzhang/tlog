package stdlog

import "sync/atomic"

type Level int32

const (
	DEBUG Level = -2
	TRACE Level = -1
	INFO  Level = 0
	WARN  Level = 1
	ERROR Level = 2
	PANIC Level = 3
	FATAL Level = 4
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case PANIC:
		return "PANIC"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

type atomicLevel struct {
	value int32
}

func (p *atomicLevel) Load() Level {
	return Level(atomic.LoadInt32(&p.value))
}

func (p *atomicLevel) Store(l Level) {
	atomic.StoreInt32(&p.value, int32(l))
}
