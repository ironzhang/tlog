package stdlog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ironzhang/tlog/logger"
)

var ExitFunc = os.Exit

var PanicFunc = func(v interface{}) {
	panic(v)
}

type field struct {
	key   string
	value interface{}
}

type Logger struct {
	base          *log.Logger
	level         *atomicLevel
	hook          logger.ContextHookFunc
	calldepth     int
	argsFields    []field
	contextFields []field
}

func NewLogger(base *log.Logger, opts ...Option) *Logger {
	if base == nil {
		panic("base is nil")
	}
	l := &Logger{
		base:      base,
		level:     newAtomicLevel(INFO),
		hook:      nil,
		calldepth: 0,
	}
	for _, fn := range opts {
		fn(l)
	}
	return l
}

func (p *Logger) SetLogger(base *log.Logger) {
	if base == nil {
		panic("base is nil")
	}
	p.base = base
}

func (p *Logger) GetLevel() Level {
	return p.level.Load()
}

func (p *Logger) SetLevel(l Level) {
	p.level.Store(l)
}

func (p *Logger) SetCalldepth(calldepth int) {
	p.calldepth = calldepth
}

func (p *Logger) GetCalldepth() int {
	return p.calldepth
}

func (p *Logger) SetContextHook(hook logger.ContextHookFunc) {
	p.hook = hook
}

func (p *Logger) clone() *Logger {
	c := &Logger{
		base:          p.base,
		level:         p.level,
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
	if p.hook == nil {
		return p
	}
	args := p.hook(ctx)
	if len(args) <= 0 {
		return p
	}

	c := p.clone()
	c.contextFields = append(c.contextFields, sweetenFields(args)...)
	return c
}

func (p *Logger) Debug(args ...interface{}) {
	p.Output(DEBUG, 2, args...)
}

func (p *Logger) Debugf(format string, args ...interface{}) {
	p.Outputf(DEBUG, 2, format, args...)
}

func (p *Logger) Debugw(message string, kvs ...interface{}) {
	p.Outputw(DEBUG, 2, message, kvs...)
}

func (p *Logger) Trace(args ...interface{}) {
	p.Output(TRACE, 2, args...)
}

func (p *Logger) Tracef(format string, args ...interface{}) {
	p.Outputf(TRACE, 2, format, args...)
}

func (p *Logger) Tracew(message string, kvs ...interface{}) {
	p.Outputw(TRACE, 2, message, kvs...)
}

func (p *Logger) Info(args ...interface{}) {
	p.Output(INFO, 2, args...)
}

func (p *Logger) Infof(format string, args ...interface{}) {
	p.Outputf(INFO, 2, format, args...)
}

func (p *Logger) Infow(message string, kvs ...interface{}) {
	p.Outputw(INFO, 2, message, kvs...)
}

func (p *Logger) Warn(args ...interface{}) {
	p.Output(WARN, 2, args...)
}

func (p *Logger) Warnf(format string, args ...interface{}) {
	p.Outputf(WARN, 2, format, args...)
}

func (p *Logger) Warnw(message string, kvs ...interface{}) {
	p.Outputw(WARN, 2, message, kvs...)
}

func (p *Logger) Error(args ...interface{}) {
	p.Output(ERROR, 2, args...)
}

func (p *Logger) Errorf(format string, args ...interface{}) {
	p.Outputf(ERROR, 2, format, args...)
}

func (p *Logger) Errorw(message string, kvs ...interface{}) {
	p.Outputw(ERROR, 2, message, kvs...)
}

func (p *Logger) Panic(args ...interface{}) {
	p.Output(PANIC, 2, args...)
}

func (p *Logger) Panicf(format string, args ...interface{}) {
	p.Outputf(PANIC, 2, format, args...)
}

func (p *Logger) Panicw(message string, kvs ...interface{}) {
	p.Outputw(PANIC, 2, message, kvs...)
}

func (p *Logger) Fatal(args ...interface{}) {
	p.Output(FATAL, 2, args...)
}

func (p *Logger) Fatalf(format string, args ...interface{}) {
	p.Outputf(FATAL, 2, format, args...)
}

func (p *Logger) Fatalw(message string, kvs ...interface{}) {
	p.Outputw(FATAL, 2, message, kvs...)
}

func (p *Logger) Output(lv Level, calldepth int, args ...interface{}) error {
	if p.level.Load() <= lv {
		d := p.calldepth + calldepth + 1
		s := p.sprint(lv, args...)
		return p.base.Output(d, s)
	}
	return nil
}

func (p *Logger) Outputf(lv Level, calldepth int, format string, args ...interface{}) error {
	if p.level.Load() <= lv {
		d := p.calldepth + calldepth + 1
		s := p.sprintf(lv, format, args...)
		return p.base.Output(d, s)
	}
	return nil
}

func (p *Logger) Outputw(lv Level, calldepth int, message string, kvs ...interface{}) error {
	if p.level.Load() <= lv {
		d := p.calldepth + calldepth + 1
		s := p.sprintw(lv, message, kvs...)
		return p.base.Output(d, s)
	}
	return nil
}

func (p *Logger) sprint(lv Level, args ...interface{}) string {
	var buf bytes.Buffer
	buf.WriteString("[" + lv.String() + "] ")
	fmt.Fprint(&buf, args...)
	buf.WriteByte('\t')
	p.writeFields(&buf)
	return buf.String()
}

func (p *Logger) sprintf(lv Level, format string, args ...interface{}) string {
	var buf bytes.Buffer
	buf.WriteString("[" + lv.String() + "] ")
	fmt.Fprintf(&buf, format, args...)
	buf.WriteByte('\t')
	p.writeFields(&buf)
	return buf.String()
}

func (p *Logger) sprintw(lv Level, message string, kvs ...interface{}) string {
	var buf bytes.Buffer
	buf.WriteString("[" + lv.String() + "] ")
	buf.WriteString(message)
	buf.WriteByte('\t')
	p.writeFields(&buf, sweetenFields(kvs)...)
	return buf.String()
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

func sweetenFields(args []interface{}) []field {
	if len(args) == 0 {
		return nil
	}

	fields := make([]field, 0, len(args)/2+1)
	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			key = fmt.Sprintf("!#%d:%v#", i, args[i])
		}
		var val interface{}
		if i+1 < len(args) {
			val = args[i+1]
		} else {
			val = "!#ignored#"
		}
		fields = append(fields, field{key: key, value: val})
	}
	return fields
}
