package zsink

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

func TestParseSize(t *testing.T) {
	tests := []struct {
		s    string
		size int
		err  string
	}{
		{s: "10TB", err: "unknown unit"},
		{s: "1.0B", err: "strconv.Atoi"},
		{s: "10", size: 10},
		{s: "10B", size: 10},
		{s: "10k", size: 10 * 1024},
		{s: "10kb", size: 10 * 1024},
		{s: "2m", size: 2 * 1024 * 1024},
		{s: "2MB", size: 2 * 1024 * 1024},
		{s: "3G", size: 3 * 1024 * 1024 * 1024},
		{s: "3GB", size: 3 * 1024 * 1024 * 1024},
	}
	for i, tt := range tests {
		size, err := parseSize(tt.s)
		if !matchError(t, err, tt.err) {
			t.Errorf("%d: match error: got %v, want %v", i, err, tt.err)
			continue
		}
		if err != nil {
			t.Logf("%d: parse size: %v", i, err)
			continue
		}
		if got, want := size, tt.size; got != want {
			t.Errorf("%d: size: got %v, want %v", i, got, want)
			continue
		}
	}
}
