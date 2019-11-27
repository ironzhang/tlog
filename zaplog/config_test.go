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
