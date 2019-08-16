package stdlog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/ironzhang/tlog/logger"
)

const baseCalldepth = 2

type field struct {
	key   string
	value interface{}
}

type Logger struct {
	base          *log.Logger
	level         *atomicLevel
	calldepth     int
	argsFields    []field
	contextFields []field
}

func NewLogger(base *log.Logger, opts ...Option) *Logger {
	l := &Logger{
		base:      base,
		level:     newAtomicLevel(INFO),
		calldepth: baseCalldepth,
	}
	for _, fn := range opts {
		fn(l)
	}
	return l
}

func (p *Logger) SetLogger(l *log.Logger) {
	p.base = l
}

func (p *Logger) GetLevel() Level {
	return p.level.Load()
}

func (p *Logger) SetLevel(l Level) {
	p.level.Store(l)
}

func (p *Logger) SetCalldepth(calldepth int) {
	p.calldepth = baseCalldepth + calldepth
}

func (p *Logger) GetCalldepth() int {
	return p.calldepth - baseCalldepth
}

func (p *Logger) clone() *Logger {
	c := &Logger{
		base:          p.base,
		level:         p.level,
		calldepth:     baseCalldepth,
		argsFields:    p.argsFields,
		contextFields: p.contextFields,
	}
	return c
}

func (p *Logger) WithArgs(args ...interface{}) logger.Logger {
	if len(args) <= 0 {
		return p
	}

	c := p.clone()
	c.argsFields = append(c.argsFields, sweetenFields(args)...)
	return c
}

func (p *Logger) WithContext(ctx context.Context) logger.Logger {
	args := logger.WithContextHook(ctx)
	if len(args) <= 0 {
		return p
	}

	c := p.clone()
	c.contextFields = append(c.contextFields, sweetenFields(args)...)
	return c
}

type byteWriter interface {
	Write(p []byte) (n int, err error)
	WriteByte(c byte) error
}

func (p *Logger) writeFields(w byteWriter, fields ...field) {
	comma := false
	if len(p.contextFields) > 0 {
		w.WriteByte('{')
		writeFields(w, p.contextFields, comma)
		comma = true
	}
	if len(p.argsFields) > 0 {
		if !comma {
			w.WriteByte('{')
		}
		writeFields(w, p.argsFields, comma)
		comma = true
	}
	if len(fields) > 0 {
		if !comma {
			w.WriteByte('{')
		}
		writeFields(w, fields, comma)
		comma = true
	}
	if comma {
		w.WriteByte('}')
	}
}

func (p *Logger) print(l Level, args ...interface{}) string {
	var buf bytes.Buffer
	buf.WriteString("[" + l.String() + "] ")
	fmt.Fprint(&buf, args...)
	buf.WriteByte('\t')
	p.writeFields(&buf)
	return buf.String()
}

func (p *Logger) printf(l Level, format string, args ...interface{}) string {
	var buf bytes.Buffer
	buf.WriteString("[" + l.String() + "] ")
	fmt.Fprintf(&buf, format, args...)
	buf.WriteByte('\t')
	p.writeFields(&buf)
	return buf.String()
}

func (p *Logger) printw(l Level, message string, kvs ...interface{}) string {
	var buf bytes.Buffer
	buf.WriteString("[" + l.String() + "] ")
	buf.WriteString(message)
	buf.WriteByte('\t')
	p.writeFields(&buf, sweetenFields(kvs)...)
	return buf.String()
}

func (p *Logger) Debug(args ...interface{}) {
	if p.level.Load() <= DEBUG {
		p.base.Output(p.calldepth, p.print(DEBUG, args...))
	}
}

func (p *Logger) Debugf(format string, args ...interface{}) {
	if p.level.Load() <= DEBUG {
		p.base.Output(p.calldepth, p.printf(DEBUG, format, args...))
	}
}

func (p *Logger) Debugw(message string, kvs ...interface{}) {
	if p.level.Load() <= DEBUG {
		p.base.Output(p.calldepth, p.printw(DEBUG, message, kvs...))
	}
}

func (p *Logger) Trace(args ...interface{}) {
	if p.level.Load() <= TRACE {
		p.base.Output(p.calldepth, p.print(TRACE, args...))
	}
}

func (p *Logger) Tracef(format string, args ...interface{}) {
	if p.level.Load() <= TRACE {
		p.base.Output(p.calldepth, p.printf(TRACE, format, args...))
	}
}

func (p *Logger) Tracew(message string, kvs ...interface{}) {
	if p.level.Load() <= TRACE {
		p.base.Output(p.calldepth, p.printw(TRACE, message, kvs...))
	}
}

func (p *Logger) Info(args ...interface{}) {
	if p.level.Load() <= INFO {
		p.base.Output(p.calldepth, p.print(INFO, args...))
	}
}

func (p *Logger) Infof(format string, args ...interface{}) {
	if p.level.Load() <= INFO {
		p.base.Output(p.calldepth, p.printf(INFO, format, args...))
	}
}

func (p *Logger) Infow(message string, kvs ...interface{}) {
	if p.level.Load() <= INFO {
		p.base.Output(p.calldepth, p.printw(INFO, message, kvs...))
	}
}

func (p *Logger) Warn(args ...interface{}) {
	if p.level.Load() <= WARN {
		p.base.Output(p.calldepth, p.print(WARN, args...))
	}
}

func (p *Logger) Warnf(format string, args ...interface{}) {
	if p.level.Load() <= WARN {
		p.base.Output(p.calldepth, p.printf(WARN, format, args...))
	}
}

func (p *Logger) Warnw(message string, kvs ...interface{}) {
	if p.level.Load() <= WARN {
		p.base.Output(p.calldepth, p.printf(WARN, message, kvs...))
	}
}

func (p *Logger) Error(args ...interface{}) {
	if p.level.Load() <= ERROR {
		p.base.Output(p.calldepth, p.print(ERROR, args...))
	}

}
func (p *Logger) Errorf(format string, args ...interface{}) {
	if p.level.Load() <= ERROR {
		p.base.Output(p.calldepth, p.printf(ERROR, format, args...))
	}
}

func (p *Logger) Errorw(message string, kvs ...interface{}) {
	if p.level.Load() <= ERROR {
		p.base.Output(p.calldepth, p.printw(ERROR, message, kvs...))
	}
}

func (p *Logger) Panic(args ...interface{}) {
	if p.level.Load() <= PANIC {
		p.base.Output(p.calldepth, p.print(PANIC, args...))
	}

}
func (p *Logger) Panicf(format string, args ...interface{}) {
	if p.level.Load() <= PANIC {
		p.base.Output(p.calldepth, p.printf(PANIC, format, args...))
	}

}
func (p *Logger) Panicw(message string, kvs ...interface{}) {
	if p.level.Load() <= PANIC {
		p.base.Output(p.calldepth, p.printw(PANIC, message, kvs...))
	}
}

func (p *Logger) Fatal(args ...interface{}) {
	if p.level.Load() <= FATAL {
		p.base.Output(p.calldepth, p.print(FATAL, args...))
	}
}

func (p *Logger) Fatalf(format string, args ...interface{}) {
	if p.level.Load() <= FATAL {
		p.base.Output(p.calldepth, p.printf(FATAL, format, args...))
	}
}

func (p *Logger) Fatalw(message string, kvs ...interface{}) {
	if p.level.Load() <= FATAL {
		p.base.Output(p.calldepth, p.printw(FATAL, message, kvs...))
	}
}

func sweetenFields(args []interface{}) []field {
	if len(args) == 0 {
		return nil
	}

	fields := make([]field, 0, len(args)/2+1)
	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			key = fmt.Sprintf("!{%d:%v}", i, args[i])
		}
		var val interface{}
		if i+1 < len(args) {
			val = args[i+1]
		} else {
			val = "!{ignored}"
		}
		fields = append(fields, field{key: key, value: val})
	}
	return fields
}

func writeFields(w io.Writer, fields []field, comma bool) {
	for _, f := range fields {
		if comma {
			fmt.Fprintf(w, ", ")
		} else {
			comma = true
		}
		fmt.Fprintf(w, "%q: %s", f.key, marshal(f.value))
	}
}

func marshal(a interface{}) string {
	switch v := a.(type) {
	case error:
		return "\"" + v.Error() + "\""
	case fmt.Stringer:
		return "\"" + v.String() + "\""
	default:
		data, _ := json.Marshal(a)
		return string(data)
	}
}
