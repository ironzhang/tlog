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
		{s: "", l: DisableStacktrace},
		{s: "disable", l: DisableStacktrace},
		{s: "warn", l: WarnStacktrace},
		{s: "error", l: ErrorStacktrace},
		{s: "ERROR", l: ErrorStacktrace},
		{s: "Error", l: ErrorStacktrace},
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
}

func TestCallerEncoderUnmarshal(t *testing.T) {
}

func TestNameEncoderMarshal(t *testing.T) {
}

func TestNameEncoderUnmarshal(t *testing.T) {
}
