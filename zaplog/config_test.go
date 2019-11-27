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

func TestStacktraceLevel(t *testing.T) {
	tests := []struct {
		l   StacktraceLevel
		s   string
		err string
	}{
		{
			l:   -1,
			s:   "StacktraceLevel(-1)",
			err: "unrecognized stacktrace level",
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
		s := tt.l.String()
		if got, want := s, tt.s; got != want {
			t.Errorf("%d: string: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: string: got %s", i, s)

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

		var l StacktraceLevel
		err = l.UnmarshalText([]byte(tt.s))
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
