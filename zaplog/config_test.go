package zaplog

import (
	"regexp"
	"testing"
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
		{
			l: -1,
			s: "StacktraceLevel(-1)",
		},
		{
			l: DisableStacktrace,
			s: "disable",
		},
		{
			l: WarnStacktrace,
			s: "warn",
		},
		{
			l: ErrorStacktrace,
			s: "error",
		},
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
		{
			s:   "unknown",
			err: "unrecognized stacktrace level",
		},
		{
			s:   "StacktraceLevel(-1)",
			err: "unrecognized stacktrace level",
		},
		{
			s: "disable",
			l: DisableStacktrace,
		},
		{
			s: "warn",
			l: WarnStacktrace,
		},
		{
			s: "error",
			l: ErrorStacktrace,
		},
		{
			s: "ERROR",
			l: ErrorStacktrace,
		},
		{
			s: "Error",
			l: ErrorStacktrace,
		},
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

/*
func TestLevelEncoderZap(t *testing.T) {
	tests := []struct {
		e LevelEncoder
		z zapcore.LevelEncoder
	}{
		{
			e: CapitalLevelEncoder,
			z: zapcore.CapitalLevelEncoder,
		},
	}
	for i, tt := range tests {
		z := tt.e.zap()
		e := zapcore.NewMapObjectEncoder()
		z(zapcore.DebugLevel, e)
		if got, want := z, tt.z; !assert.Equal(t, got, want) {
			t.Errorf("%d: zap: got %v, want %v", i, got, want)
			continue
		}
	}
}
*/

func TestLevelEncoderMarshal(t *testing.T) {
	tests := []struct {
		e LevelEncoder
		s string
	}{
		{-1, "LevelEncoder(-1)"},
		{CapitalLevelEncoder, "capital"},
		{CapitalColorLevelEncoder, "capitalColor"},
		{LowercaseLevelEncoder, "lowercase"},
		{LowercaseColorLevelEncoder, "lowercaseColor"},
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
