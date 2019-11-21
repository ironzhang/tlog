package zaplog

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/ironzhang/tlog/iface"
	"go.uber.org/zap/zapcore"
)

var (
	errUnmarshalNilStacktraceLevel = errors.New("can't unmarshal a nil *StacktraceLevel")
	errUnmarshalNilLevelEncoder    = errors.New("can't unmarshal a nil *LevelEncoder")
	errUnmarshalNilTimeEncoder     = errors.New("can't unmarshal a nil *TimeEncoder")
	errUnmarshalNilDurationEncoder = errors.New("can't unmarshal a nil *DurationEncoder")
	errUnmarshalNilCallerEncoder   = errors.New("can't unmarshal a nil *CallerEncoder")
	errUnmarshalNilNameEncoder     = errors.New("can't unmarshal a nil *NameEncoder")
)

type StacktraceLevel int8

const (
	DisableStacktrace StacktraceLevel = iota
	WarnStacktrace
	ErrorStacktrace
)

func (l StacktraceLevel) String() string {
	switch l {
	case DisableStacktrace:
		return "disable"
	case WarnStacktrace:
		return "warn"
	case ErrorStacktrace:
		return "error"
	default:
		return fmt.Sprintf("StacktraceLevel(%d)", l)
	}
}

func (l StacktraceLevel) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

func (l *StacktraceLevel) UnmarshalText(text []byte) error {
	if l == nil {
		return errUnmarshalNilStacktraceLevel
	}
	if !l.unmarshalText(text) && !l.unmarshalText(bytes.ToLower(text)) {
		return fmt.Errorf("unrecognized stacktrace level %q", text)
	}
	return nil
}

func (l *StacktraceLevel) unmarshalText(text []byte) bool {
	switch string(text) {
	case "disable", "DISABLE", "":
		*l = DisableStacktrace
	case "warn", "WARN":
		*l = WarnStacktrace
	case "error", "ERROR":
		*l = ErrorStacktrace
	default:
		return false
	}
	return true
}

type LevelEncoder int8

const (
	CapitalLevelEncoder LevelEncoder = iota
	CapitalColorLevelEncoder
	LowercaseLevelEncoder
	LowercaseColorLevelEncoder
)

func (e LevelEncoder) zap() zapcore.LevelEncoder {
	switch e {
	case CapitalLevelEncoder:
		return zapcore.CapitalLevelEncoder
	case CapitalColorLevelEncoder:
		return zapcore.CapitalColorLevelEncoder
	case LowercaseLevelEncoder:
		return zapcore.LowercaseLevelEncoder
	case LowercaseColorLevelEncoder:
		return zapcore.LowercaseColorLevelEncoder
	default:
		return zapcore.CapitalLevelEncoder
	}
}

func (e LevelEncoder) String() string {
	switch e {
	case CapitalLevelEncoder:
		return "capital"
	case CapitalColorLevelEncoder:
		return "capitalColor"
	case LowercaseLevelEncoder:
		return "lowercase"
	case LowercaseColorLevelEncoder:
		return "lowercaseColor"
	default:
		return fmt.Sprintf("LevelEncoder(%d)", e)
	}
}

func (e LevelEncoder) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

func (e *LevelEncoder) UnmarshalText(text []byte) error {
	if e == nil {
		return errUnmarshalNilLevelEncoder
	}
	if !e.unmarshalText(text) && !e.unmarshalText(bytes.ToLower(text)) {
		return fmt.Errorf("unrecognized level encoder %q", text)
	}
	return nil
}

func (e *LevelEncoder) unmarshalText(text []byte) bool {
	switch string(text) {
	case "capital", "CAPITAL", "":
		*e = CapitalLevelEncoder
	case "capitalColor", "capitalcolor", "CAPITALCOLOR":
		*e = CapitalColorLevelEncoder
	case "lowercase", "LOWERCASE":
		*e = LowercaseLevelEncoder
	case "lowercaseColor", "lowercasecolor", "LOWERCASECOLOR":
		*e = LowercaseColorLevelEncoder
	default:
		return false
	}
	return true
}

type TimeEncoder int8

const (
	ISO8601TimeEncoder TimeEncoder = iota
	EpochTimeEncoder
	EpochNanosTimeEncoder
	EpochMillisTimeEncoder
	RFC3339TimeEncoder
	RFC3339NanoTimeEncoder
)

func (e TimeEncoder) zap() zapcore.TimeEncoder {
	switch e {
	case ISO8601TimeEncoder:
		return zapcore.ISO8601TimeEncoder
	case EpochTimeEncoder:
		return zapcore.EpochTimeEncoder
	case EpochNanosTimeEncoder:
		return zapcore.EpochNanosTimeEncoder
	case EpochMillisTimeEncoder:
		return zapcore.EpochMillisTimeEncoder
	case RFC3339TimeEncoder:
		return zapcore.RFC3339TimeEncoder
	case RFC3339NanoTimeEncoder:
		return zapcore.RFC3339NanoTimeEncoder
	default:
		return zapcore.ISO8601TimeEncoder
	}
}

func (e TimeEncoder) String() string {
	switch e {
	case ISO8601TimeEncoder:
		return "iso8601"
	case EpochTimeEncoder:
		return "epoch"
	case EpochNanosTimeEncoder:
		return "epochNanos"
	case EpochMillisTimeEncoder:
		return "epochMillis"
	case RFC3339TimeEncoder:
		return "rfc3339"
	case RFC3339NanoTimeEncoder:
		return "rfc3339nano"
	default:
		return fmt.Sprintf("TimeEncoder(%d)", e)
	}
}

func (e TimeEncoder) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

func (e *TimeEncoder) UnmarshalText(text []byte) error {
	if e == nil {
		return errUnmarshalNilTimeEncoder
	}
	if !e.unmarshalText(text) && !e.unmarshalText(text) {
		return fmt.Errorf("unrecognized time encoder %q", text)
	}
	return nil
}

func (e *TimeEncoder) unmarshalText(text []byte) bool {
	switch string(text) {
	case "iso8601", "ISO8601", "":
		*e = ISO8601TimeEncoder
	case "epoch", "EPOCH":
		*e = EpochTimeEncoder
	case "epochNanos", "epochnanos", "EPOCHNANOS":
		*e = EpochNanosTimeEncoder
	case "epochMillis", "epochmillis", "EPOCHMILLIS":
		*e = EpochMillisTimeEncoder
	case "rfc3339", "RFC3339":
		*e = RFC3339TimeEncoder
	case "rfc3339nano", "RFC3339NANO":
		*e = RFC3339NanoTimeEncoder
	default:
		return false
	}
	return true
}

type DurationEncoder int8

const (
	StringDurationEncoder DurationEncoder = iota
	SecondsDurationEncoder
	NanosDurationEncoder
)

func (e DurationEncoder) zap() zapcore.DurationEncoder {
	switch e {
	case StringDurationEncoder:
		return zapcore.StringDurationEncoder
	case SecondsDurationEncoder:
		return zapcore.SecondsDurationEncoder
	case NanosDurationEncoder:
		return zapcore.NanosDurationEncoder
	default:
		return zapcore.StringDurationEncoder
	}
}

func (e DurationEncoder) String() string {
	switch e {
	case StringDurationEncoder:
		return "string"
	case SecondsDurationEncoder:
		return "seconds"
	case NanosDurationEncoder:
		return "nanos"
	default:
		return fmt.Sprintf("DurationEncoder(%d)", e)
	}
}

func (e DurationEncoder) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

func (e *DurationEncoder) UnmarshalText(text []byte) error {
	if e == nil {
		return errUnmarshalNilDurationEncoder
	}
	if !e.unmarshalText(text) && !e.unmarshalText(text) {
		return fmt.Errorf("unrecognized duration encoder %q", text)
	}
	return nil
}

func (e *DurationEncoder) unmarshalText(text []byte) bool {
	switch string(text) {
	case "string", "STRING", "":
		*e = StringDurationEncoder
	case "seconds", "SECONDS":
		*e = SecondsDurationEncoder
	case "nanos", "NANOS":
		*e = NanosDurationEncoder
	default:
		return false
	}
	return true
}

type CallerEncoder int8

const (
	ShortCallerEncoder CallerEncoder = iota
	FullCallerEncoder
)

func (e CallerEncoder) zap() zapcore.CallerEncoder {
	switch e {
	case ShortCallerEncoder:
		return zapcore.ShortCallerEncoder
	case FullCallerEncoder:
		return zapcore.FullCallerEncoder
	default:
		return zapcore.ShortCallerEncoder
	}
}

func (e CallerEncoder) String() string {
	switch e {
	case ShortCallerEncoder:
		return "short"
	case FullCallerEncoder:
		return "full"
	default:
		return fmt.Sprintf("CallerEncoder(%d)", e)
	}
}

func (e CallerEncoder) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

func (e *CallerEncoder) UnmarshalText(text []byte) error {
	if e == nil {
		return errUnmarshalNilCallerEncoder
	}
	if !e.unmarshalText(text) && !e.unmarshalText(text) {
		return fmt.Errorf("unrecognized caller encoder %q", text)
	}
	return nil
}

func (e *CallerEncoder) unmarshalText(text []byte) bool {
	switch string(text) {
	case "short", "SHORT", "":
		*e = ShortCallerEncoder
	case "full", "FULL":
		*e = FullCallerEncoder
	default:
		return false
	}
	return true
}

type NameEncoder int8

const (
	FullNameEncoder NameEncoder = iota
)

func (e NameEncoder) zap() zapcore.NameEncoder {
	switch e {
	case FullNameEncoder:
		return zapcore.FullNameEncoder
	default:
		return zapcore.FullNameEncoder
	}
}

func (e NameEncoder) String() string {
	switch e {
	case FullNameEncoder:
		return "full"
	default:
		return fmt.Sprintf("NameEncoder(%d)", e)
	}
}

func (e NameEncoder) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

func (e *NameEncoder) UnmarshalText(text []byte) error {
	if e == nil {
		return errUnmarshalNilNameEncoder
	}
	if !e.unmarshalText(text) && !e.unmarshalText(text) {
		return fmt.Errorf("unrecognized name encoder %q", text)
	}
	return nil
}

func (e *NameEncoder) unmarshalText(text []byte) bool {
	switch string(text) {
	case "full", "FULL", "":
		*e = FullNameEncoder
	default:
		return false
	}
	return true
}

type SinkConfig struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url" yaml:"url"`
}

type EncoderConfig struct {
	MessageKey     string          `json:"messageKey" yaml:"messageKey"`
	LevelKey       string          `json:"levelKey" yaml:"levelKey"`
	TimeKey        string          `json:"timeKey" yaml:"timeKey"`
	NameKey        string          `json:"nameKey" yaml:"nameKey"`
	CallerKey      string          `json:"callerKey" yaml:"callerKey"`
	StacktraceKey  string          `json:"stacktraceKey" yaml:"stacktraceKey"`
	EncodeLevel    LevelEncoder    `json:"levelEncoder" yaml:"levelEncoder"`
	EncodeTime     TimeEncoder     `json:"timeEncoder" yaml:"timeEncoder"`
	EncodeDuration DurationEncoder `json:"durationEncoder" yaml:"durationEncoder"`
	EncodeCaller   CallerEncoder   `json:"callerEncoder" yaml:"callerEncoder"`
	EncodeName     NameEncoder     `json:"nameEncoder" yaml:"nameEncoder"`
}

type CoreConfig struct {
	Name     string        `json:"name" yaml:"name"`
	Encoding string        `json:"encoding" yaml:"encoding"`
	Encoder  EncoderConfig `json:"encoder" yaml:"encoder"`
	MinLevel iface.Level   `json:"minLevel" yaml:"minLevel"`
	MaxLevel iface.Level   `json:"maxLevel" yaml:"maxLevel"`
	Sinks    []string      `json:"sinks" yaml:"sinks"`
}

type LoggerConfig struct {
	Name            string          `json:"name" yaml:"name"`
	DisableCaller   bool            `json:"disableCaller" yaml:"disableCaller"`
	StacktraceLevel StacktraceLevel `json:"stacktraceLevel" yaml:"stacktraceLevel"`
	Cores           []string        `json:"cores" yaml:"cores"`
}

type Config struct {
	Level   iface.Level    `json:"level" yaml:"level"`
	Sinks   []SinkConfig   `json:"sinks" yaml:"sinks"`
	Cores   []CoreConfig   `json:"cores" yaml:"cores"`
	Loggers []LoggerConfig `json:"loggers" yaml:"loggers"`
}
