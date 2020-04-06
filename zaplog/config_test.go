package zaplog

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func matchError(t testing.TB, err error, errstr string) bool {
	switch {
	case errstr == "" && err == nil:
		return true
	case errstr == "" && err != nil:
		return false
	case errstr != "" && err == nil:
		return false
	case errstr != "" && err != nil:
		matched, err := regexp.MatchString(errstr, err.Error())
		if err != nil {
			t.Fatalf("match error: %v", err)
		}
		return matched
	}
	return false
}

func TestStacktraceLevelMarshal(t *testing.T) {
	tests := []struct {
		l StacktraceLevel
		s string
	}{
		{l: -1, s: "StacktraceLevel(-1)"},
		{l: DisableStacktrace, s: "disable"},
		{l: WarnStacktrace, s: "warn"},
		{l: ErrorStacktrace, s: "error"},
	}
	for i, tt := range tests {
		text, err := tt.l.MarshalText()
		if err != nil {
			t.Errorf("%d: marshal text: %v", i, err)
			continue
		}
		if got, want := string(text), tt.s; got != want {
			t.Errorf("%d: text: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: text: got %s", i, text)
	}
}

func TestStacktraceLevelUnmarshal(t *testing.T) {
	tests := []struct {
		s   string
		l   StacktraceLevel
		err string
	}{
		{s: "unknown", err: "unrecognized stacktrace level"},
		{s: "StacktraceLevel(-1)", err: "unrecognized stacktrace level"},
		{s: "", l: PanicStacktrace},
		{s: "panic", l: PanicStacktrace},
		{s: "error", l: ErrorStacktrace},
		{s: "warn", l: WarnStacktrace},
		{s: "disable", l: DisableStacktrace},
		{s: "DISABLE", l: DisableStacktrace},
		{s: "dISABLE", l: DisableStacktrace},
	}
	for i, tt := range tests {
		var l StacktraceLevel
		err := l.UnmarshalText([]byte(tt.s))
		if !matchError(t, err, tt.err) {
			t.Errorf("%d: match error: got %v, want %v", i, err, tt.err)
			continue
		}
		if err != nil {
			t.Logf("%d: unmarshal text: %v", i, err)
			continue
		}
		if got, want := l, tt.l; got != want {
			t.Errorf("%d: stacktrace level: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: stacktrace level: got %v", i, l)
	}
}

func TestLevelEncoderMarshal(t *testing.T) {
	tests := []struct {
		e LevelEncoder
		s string
	}{
		{e: -1, s: "LevelEncoder(-1)"},
		{e: CapitalLevelEncoder, s: "capital"},
		{e: CapitalColorLevelEncoder, s: "capitalColor"},
		{e: LowercaseLevelEncoder, s: "lowercase"},
		{e: LowercaseColorLevelEncoder, s: "lowercaseColor"},
	}
	for i, tt := range tests {
		text, err := tt.e.MarshalText()
		if err != nil {
			t.Errorf("%d: marshal text: %v", i, err)
			continue
		}
		if got, want := string(text), tt.s; got != want {
			t.Errorf("%d: text: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: text: got %s", i, text)
	}
}

func TestLevelEncoderUnmarshal(t *testing.T) {
	tests := []struct {
		s   string
		e   LevelEncoder
		err string
	}{
		{s: "LevelEncoder(-1)", err: "unrecognized level encoder"},
		{s: "", e: CapitalLevelEncoder},
		{s: "capital", e: CapitalLevelEncoder},
		{s: "capitalColor", e: CapitalColorLevelEncoder},
		{s: "lowercase", e: LowercaseLevelEncoder},
		{s: "lowercaseColor", e: LowercaseColorLevelEncoder},
		{s: "LOWERCASECOLOR", e: LowercaseColorLevelEncoder},
		{s: "lowercasecolor", e: LowercaseColorLevelEncoder},
		{s: "loWercasecolor", e: LowercaseColorLevelEncoder},
	}
	for i, tt := range tests {
		var e LevelEncoder
		err := e.UnmarshalText([]byte(tt.s))
		if !matchError(t, err, tt.err) {
			t.Errorf("%d: match error: got %v, want %v", i, err, tt.err)
			continue
		}
		if err != nil {
			t.Logf("%d: unmarshal text: %v", i, err)
			continue
		}
		if got, want := e, tt.e; got != want {
			t.Errorf("%d: level encoder: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: level encoder: got %v", i, e)
	}
}

func TestTimeEncoderMarshal(t *testing.T) {
	tests := []struct {
		e TimeEncoder
		s string
	}{
		{e: -1, s: "TimeEncoder(-1)"},
		{e: ISO8601TimeEncoder, s: "iso8601"},
		{e: EpochTimeEncoder, s: "epoch"},
		{e: EpochNanosTimeEncoder, s: "epochNanos"},
		{e: EpochMillisTimeEncoder, s: "epochMillis"},
		{e: RFC3339TimeEncoder, s: "rfc3339"},
		{e: RFC3339NanoTimeEncoder, s: "rfc3339nano"},
	}
	for i, tt := range tests {
		text, err := tt.e.MarshalText()
		if err != nil {
			t.Errorf("%d: marshal text: %v", i, err)
			continue
		}
		if got, want := string(text), tt.s; got != want {
			t.Errorf("%d: text: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: text: got %s", i, text)
	}
}

func TestTimeEncoderUnmarshal(t *testing.T) {
	tests := []struct {
		s   string
		e   TimeEncoder
		err string
	}{
		{s: "TimeEncoder(-1)", err: "unrecognized time encoder"},
		{s: "", e: ISO8601TimeEncoder},
		{s: "iso8601", e: ISO8601TimeEncoder},
		{s: "epoch", e: EpochTimeEncoder},
		{s: "epochNanos", e: EpochNanosTimeEncoder},
		{s: "epochMillis", e: EpochMillisTimeEncoder},
		{s: "rfc3339", e: RFC3339TimeEncoder},
		{s: "rfc3339nano", e: RFC3339NanoTimeEncoder},
		{s: "RFC3339NANO", e: RFC3339NanoTimeEncoder},
		{s: "rFC3339NANO", e: RFC3339NanoTimeEncoder},
	}
	for i, tt := range tests {
		var e TimeEncoder
		err := e.UnmarshalText([]byte(tt.s))
		if !matchError(t, err, tt.err) {
			t.Errorf("%d: match error: got %v, want %v", i, err, tt.err)
			continue
		}
		if err != nil {
			t.Logf("%d: unmarshal text: %v", i, err)
			continue
		}
		if got, want := e, tt.e; got != want {
			t.Errorf("%d: time encoder: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: time encoder: got %v", i, e)
	}
}

func TestDurationEncoderMarshal(t *testing.T) {
	tests := []struct {
		e DurationEncoder
		s string
	}{
		{e: -1, s: "DurationEncoder(-1)"},
		{e: StringDurationEncoder, s: "string"},
		{e: SecondsDurationEncoder, s: "seconds"},
		{e: NanosDurationEncoder, s: "nanos"},
	}
	for i, tt := range tests {
		text, err := tt.e.MarshalText()
		if err != nil {
			t.Errorf("%d: marshal text: %v", i, err)
			continue
		}
		if got, want := string(text), tt.s; got != want {
			t.Errorf("%d: text: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: text: got %s", i, text)
	}
}

func TestDurationEncoderUnmarshal(t *testing.T) {
	tests := []struct {
		s   string
		e   DurationEncoder
		err string
	}{
		{s: "DurationEncoder(-1)", err: "unrecognized duration encoder"},
		{s: "", e: StringDurationEncoder},
		{s: "string", e: StringDurationEncoder},
		{s: "seconds", e: SecondsDurationEncoder},
		{s: "nanos", e: NanosDurationEncoder},
		{s: "NANOS", e: NanosDurationEncoder},
		{s: "nANOS", e: NanosDurationEncoder},
	}
	for i, tt := range tests {
		var e DurationEncoder
		err := e.UnmarshalText([]byte(tt.s))
		if !matchError(t, err, tt.err) {
			t.Errorf("%d: match error: got %v, want %v", i, err, tt.err)
			continue
		}
		if err != nil {
			t.Logf("%d: unmarshal text: %v", i, err)
			continue
		}
		if got, want := e, tt.e; got != want {
			t.Errorf("%d: duration encoder: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: duration encoder: got %v", i, e)
	}
}

func TestCallerEncoderMarshal(t *testing.T) {
	tests := []struct {
		e CallerEncoder
		s string
	}{
		{e: -1, s: "CallerEncoder(-1)"},
		{e: ShortCallerEncoder, s: "short"},
		{e: FullCallerEncoder, s: "full"},
	}
	for i, tt := range tests {
		text, err := tt.e.MarshalText()
		if err != nil {
			t.Errorf("%d: marshal text: %v", i, err)
			continue
		}
		if got, want := string(text), tt.s; got != want {
			t.Errorf("%d: text: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: text: got %s", i, text)
	}
}

func TestCallerEncoderUnmarshal(t *testing.T) {
	tests := []struct {
		s   string
		e   CallerEncoder
		err string
	}{
		{s: "CallerEncoder(-1)", err: "unrecognized caller encoder"},
		{s: "", e: ShortCallerEncoder},
		{s: "short", e: ShortCallerEncoder},
		{s: "full", e: FullCallerEncoder},
		{s: "FULL", e: FullCallerEncoder},
		{s: "fULL", e: FullCallerEncoder},
	}
	for i, tt := range tests {
		var e CallerEncoder
		err := e.UnmarshalText([]byte(tt.s))
		if !matchError(t, err, tt.err) {
			t.Errorf("%d: match error: got %v, want %v", i, err, tt.err)
			continue
		}
		if err != nil {
			t.Logf("%d: unmarshal text: %v", i, err)
			continue
		}
		if got, want := e, tt.e; got != want {
			t.Errorf("%d: caller encoder: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: caller encoder: got %v", i, e)
	}
}

func TestNameEncoderMarshal(t *testing.T) {
	tests := []struct {
		e NameEncoder
		s string
	}{
		{e: -1, s: "NameEncoder(-1)"},
		{e: FullNameEncoder, s: "full"},
	}
	for i, tt := range tests {
		text, err := tt.e.MarshalText()
		if err != nil {
			t.Errorf("%d: marshal text: %v", i, err)
			continue
		}
		if got, want := string(text), tt.s; got != want {
			t.Errorf("%d: text: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: text: got %s", i, text)
	}
}

func TestNameEncoderUnmarshal(t *testing.T) {
	tests := []struct {
		s   string
		e   NameEncoder
		err string
	}{
		{s: "NameEncoder(-1)", err: "unrecognized name encoder"},
		{s: "", e: FullNameEncoder},
		{s: "full", e: FullNameEncoder},
		{s: "Full", e: FullNameEncoder},
		{s: "fULL", e: FullNameEncoder},
	}
	for i, tt := range tests {
		var e NameEncoder
		err := e.UnmarshalText([]byte(tt.s))
		if !matchError(t, err, tt.err) {
			t.Errorf("%d: match error: got %v, want %v", i, err, tt.err)
			continue
		}
		if err != nil {
			t.Logf("%d: unmarshal text: %v", i, err)
			continue
		}
		if got, want := e, tt.e; got != want {
			t.Errorf("%d: name encoder: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: name encoder: got %v", i, e)
	}
}

type tPrimitiveArrayEncoder struct {
	elems []interface{}
}

func (p *tPrimitiveArrayEncoder) AppendBool(v bool)             { p.elems = append(p.elems, v) }
func (p *tPrimitiveArrayEncoder) AppendByteString(v []byte)     { p.elems = append(p.elems, v) }
func (p *tPrimitiveArrayEncoder) AppendComplex128(v complex128) { p.elems = append(p.elems, v) }
func (p *tPrimitiveArrayEncoder) AppendComplex64(v complex64)   { p.elems = append(p.elems, v) }
func (p *tPrimitiveArrayEncoder) AppendFloat64(v float64)       { p.elems = append(p.elems, v) }
func (p *tPrimitiveArrayEncoder) AppendFloat32(v float32)       { p.elems = append(p.elems, v) }
func (p *tPrimitiveArrayEncoder) AppendInt(v int)               { p.elems = append(p.elems, v) }
func (p *tPrimitiveArrayEncoder) AppendInt64(v int64)           { p.elems = append(p.elems, v) }
func (p *tPrimitiveArrayEncoder) AppendInt32(v int32)           { p.elems = append(p.elems, v) }
func (p *tPrimitiveArrayEncoder) AppendInt16(v int16)           { p.elems = append(p.elems, v) }
func (p *tPrimitiveArrayEncoder) AppendInt8(v int8)             { p.elems = append(p.elems, v) }
func (p *tPrimitiveArrayEncoder) AppendString(v string)         { p.elems = append(p.elems, v) }
func (p *tPrimitiveArrayEncoder) AppendUint(v uint)             { p.elems = append(p.elems, v) }
func (p *tPrimitiveArrayEncoder) AppendUint64(v uint64)         { p.elems = append(p.elems, v) }
func (p *tPrimitiveArrayEncoder) AppendUint32(v uint32)         { p.elems = append(p.elems, v) }
func (p *tPrimitiveArrayEncoder) AppendUint16(v uint16)         { p.elems = append(p.elems, v) }
func (p *tPrimitiveArrayEncoder) AppendUint8(v uint8)           { p.elems = append(p.elems, v) }
func (p *tPrimitiveArrayEncoder) AppendUintptr(v uintptr)       { p.elems = append(p.elems, v) }

func TestLevelEncoder(t *testing.T) {
	tests := []struct {
		encoder LevelEncoder
		elems   []interface{}
	}{
		{encoder: -1, elems: []interface{}{"INFO"}},
		{encoder: CapitalLevelEncoder, elems: []interface{}{"INFO"}},
		{encoder: LowercaseLevelEncoder, elems: []interface{}{"info"}},
		{encoder: CapitalColorLevelEncoder, elems: []interface{}{"\x1b[34mINFO\x1b[0m"}},
		{encoder: LowercaseColorLevelEncoder, elems: []interface{}{"\x1b[34minfo\x1b[0m"}},
	}
	for i, tt := range tests {
		enc := &tPrimitiveArrayEncoder{}
		tt.encoder.zap()(zapcore.InfoLevel, enc)
		assert.Equal(t, tt.elems, enc.elems, "%d: unexpected level encoder elements", i)
	}
}

func TestTimeEncoder(t *testing.T) {
	ts := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		encoder TimeEncoder
		elems   []interface{}
	}{
		{encoder: -1, elems: []interface{}{ts.Format("2006-01-02T15:04:05.000Z0700")}},
		{encoder: ISO8601TimeEncoder, elems: []interface{}{ts.Format("2006-01-02T15:04:05.000Z0700")}},
		{encoder: EpochTimeEncoder, elems: []interface{}{float64(0)}},
		{encoder: EpochNanosTimeEncoder, elems: []interface{}{int64(0)}},
		{encoder: EpochMillisTimeEncoder, elems: []interface{}{float64(0)}},
		{encoder: RFC3339TimeEncoder, elems: []interface{}{ts.Format(time.RFC3339)}},
		{encoder: RFC3339NanoTimeEncoder, elems: []interface{}{ts.Format(time.RFC3339Nano)}},
	}
	for i, tt := range tests {
		enc := &tPrimitiveArrayEncoder{}
		tt.encoder.zap()(ts, enc)
		assert.Equal(t, tt.elems, enc.elems, "%d: unexpected time encoder elements", i)
	}
}

func TestDurationEncoder(t *testing.T) {
	d := time.Minute
	tests := []struct {
		encoder DurationEncoder
		elems   []interface{}
	}{
		{encoder: -1, elems: []interface{}{d.String()}},
		{encoder: StringDurationEncoder, elems: []interface{}{d.String()}},
		{encoder: SecondsDurationEncoder, elems: []interface{}{float64(d) / float64(time.Second)}},
		{encoder: NanosDurationEncoder, elems: []interface{}{int64(d)}},
	}
	for i, tt := range tests {
		enc := &tPrimitiveArrayEncoder{}
		tt.encoder.zap()(d, enc)
		assert.Equal(t, tt.elems, enc.elems, "%d: unexpected duration encoder elements", i)
	}
}

func TestCallerEncoder(t *testing.T) {
	caller := zapcore.EntryCaller{
		Defined: true,
		File:    "git.xiaojukeji.com/pearls/tlog/zaplog/config_test.go",
		Line:    460,
	}
	tests := []struct {
		encoder CallerEncoder
		elems   []interface{}
	}{
		{encoder: -1, elems: []interface{}{caller.TrimmedPath()}},
		{encoder: ShortCallerEncoder, elems: []interface{}{caller.TrimmedPath()}},
		{encoder: FullCallerEncoder, elems: []interface{}{caller.String()}},
	}
	for i, tt := range tests {
		enc := &tPrimitiveArrayEncoder{}
		tt.encoder.zap()(caller, enc)
		assert.Equal(t, tt.elems, enc.elems, "%d: unexpected caller encoder elements", i)
	}
}

func TestNameEncoder(t *testing.T) {
	name := "access"
	tests := []struct {
		encoder NameEncoder
		elems   []interface{}
	}{
		{encoder: -1, elems: []interface{}{name}},
		{encoder: FullNameEncoder, elems: []interface{}{name}},
	}
	for i, tt := range tests {
		enc := &tPrimitiveArrayEncoder{}
		tt.encoder.zap()(name, enc)
		assert.Equal(t, tt.elems, enc.elems, "%d: unexpected name encoder elements", i)
	}
}
