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
	PanicStacktrace StacktraceLevel = iota
	ErrorStacktrace
	WarnStacktrace
	DisableStacktrace
)

func (l StacktraceLevel) String() string {
	switch l {
	case PanicStacktrace:
		return "panic"
	case ErrorStacktrace:
		return "error"
	case WarnStacktrace:
		return "warn"
	case DisableStacktrace:
		return "disable"
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
	case "panic", "PANIC", "":
		*l = PanicStacktrace
	case "error", "ERROR":
		*l = ErrorStacktrace
	case "warn", "WARN":
		*l = WarnStacktrace
	case "disable", "DISABLE":
		*l = DisableStacktrace
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
	if !e.unmarshalText(text) && !e.unmarshalText(bytes.ToLower(text)) {
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
	if !e.unmarshalText(text) && !e.unmarshalText(bytes.ToLower(text)) {
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
	if !e.unmarshalText(text) && !e.unmarshalText(bytes.ToLower(text)) {
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
	if !e.unmarshalText(text) && !e.unmarshalText(bytes.ToLower(text)) {
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

type EncoderConfig struct {
	MessageKey     string          `json:"messageKey,omitempty" yaml:"messageKey,omitempty"`
	LevelKey       string          `json:"levelKey,omitempty" yaml:"levelKey,omitempty"`
	TimeKey        string          `json:"timeKey,omitempty" yaml:"timeKey,omitempty"`
	NameKey        string          `json:"nameKey,omitempty" yaml:"nameKey,omitempty"`
	CallerKey      string          `json:"callerKey,omitempty" yaml:"callerKey,omitempty"`
	StacktraceKey  string          `json:"stacktraceKey,omitempty" yaml:"stacktraceKey,omitempty"`
	EncodeLevel    LevelEncoder    `json:"levelEncoder,omitempty" yaml:"levelEncoder,omitempty"`
	EncodeTime     TimeEncoder     `json:"timeEncoder,omitempty" yaml:"timeEncoder,omitempty"`
	EncodeDuration DurationEncoder `json:"durationEncoder,omitempty" yaml:"durationEncoder,omitempty"`
	EncodeCaller   CallerEncoder   `json:"callerEncoder,omitempty" yaml:"callerEncoder,omitempty"`
	EncodeName     NameEncoder     `json:"nameEncoder,omitempty" yaml:"nameEncoder,omitempty"`
}

type OutputConfig struct {
	MinLevel iface.Level `json:"minLevel" yaml:"minLevel"`
	MaxLevel iface.Level `json:"maxLevel" yaml:"maxLevel"`
	URLs     []string    `json:"urls,omitempty" yaml:"urls,omitempty"`
}

type LoggerConfig struct {
	Name            string          `json:"name,omitempty" yaml:"name,omitempty"`
	DisableCaller   bool            `json:"disableCaller,omitempty" yaml:"disableCaller,omitempty"`
	StacktraceLevel StacktraceLevel `json:"stacktraceLevel,omitempty" yaml:"stacktraceLevel,omitempty"`
	Encoding        string          `json:"encoding,omitempty" yaml:"encoding,omitempty"`
	Encoder         EncoderConfig   `json:"encoder,omitempty" yaml:"encoder,omitempty"`
	Outputs         []OutputConfig  `json:"outputs,omitempty" yaml:"outputs,omitempty"`
}

type Config struct {
	Level   iface.Level    `json:"level,omitempty" yaml:"level,omitempty"`
	Loggers []LoggerConfig `json:"loggers,omitempty" yaml:"loggers,omitempty"`
}

func NewConsoleEncoderConfig() EncoderConfig {
	return EncoderConfig{
		MessageKey:     "M",
		LevelKey:       "L",
		TimeKey:        "T",
		NameKey:        "N",
		CallerKey:      "C",
		StacktraceKey:  "S",
		EncodeLevel:    CapitalLevelEncoder,
		EncodeTime:     ISO8601TimeEncoder,
		EncodeDuration: StringDurationEncoder,
		EncodeCaller:   ShortCallerEncoder,
		EncodeName:     FullNameEncoder,
	}
}

func NewJSONEncoderConfig() EncoderConfig {
	return EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "ts",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    LowercaseLevelEncoder,
		EncodeTime:     EpochTimeEncoder,
		EncodeDuration: StringDurationEncoder,
		EncodeCaller:   ShortCallerEncoder,
		EncodeName:     FullNameEncoder,
	}
}

func NewDevelopmentConfig() Config {
	return Config{
		Level: iface.DEBUG,
		Loggers: []LoggerConfig{
			{
				Name:            "",
				DisableCaller:   false,
				StacktraceLevel: DisableStacktrace,
				Encoding:        "console",
				Encoder:         NewConsoleEncoderConfig(),
				Outputs: []OutputConfig{
					{
						MinLevel: iface.DEBUG,
						MaxLevel: iface.FATAL,
						URLs:     []string{"stderr"},
					},
				},
			},
		},
	}
}

func NewProductionConfig() Config {
	return Config{
		Level: iface.INFO,
		Loggers: []LoggerConfig{
			{
				Name:            "",
				DisableCaller:   false,
				StacktraceLevel: PanicStacktrace,
				Encoding:        "json",
				Encoder:         NewJSONEncoderConfig(),
				Outputs: []OutputConfig{
					{
						MinLevel: iface.DEBUG,
						MaxLevel: iface.DEBUG,
						URLs:     []string{"./log/debug.log"},
					},
					{
						MinLevel: iface.INFO,
						MaxLevel: iface.FATAL,
						URLs:     []string{"./log/info.log"},
					},
					{
						MinLevel: iface.WARN,
						MaxLevel: iface.FATAL,
						URLs:     []string{"./log/warn.log"},
					},
					{
						MinLevel: iface.ERROR,
						MaxLevel: iface.FATAL,
						URLs:     []string{"./log/error.log"},
					},
					{
						MinLevel: iface.PANIC,
						MaxLevel: iface.FATAL,
						URLs:     []string{"./log/fatal.log"},
					},
				},
			},
		},
	}
}
