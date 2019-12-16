package zsink

import (
	"regexp"
	"testing"
	"time"
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
		{s: "aa", err: "unknown unit"},
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

var TestValues = values{
	"k1": []string{"1"},
	"k2": []string{"2", "3"},
	"k3": []string{"300s"},
	"k4": []string{"3hour"},
	"k5": []string{"500m"},
	"k6": []string{"6XB"},
}

func TestValuesGet(t *testing.T) {
	tests := []struct {
		key   string
		exist bool
		val   string
	}{
		{key: "k1", exist: true, val: "1"},
		{key: "k2", exist: true, val: "2"},
		{key: "nokey", exist: false},
	}
	for i, tt := range tests {
		val, exist := TestValues.Get(tt.key)
		if exist != tt.exist {
			t.Errorf("%d: exist: got %v, want %v", i, exist, tt.exist)
			continue
		}
		if exist {
			if got, want := val, tt.val; got != want {
				t.Errorf("%d: value: got %v, want %v", i, got, want)
				continue
			}
		}
	}
}

func TestValuesGetDuration(t *testing.T) {
	tests := []struct {
		key   string
		err   string
		exist bool
		d     time.Duration
	}{
		{key: "nokey", exist: false},
		{key: "k3", exist: true, d: 300 * time.Second},
		{key: "k4", err: "unknown unit"},
	}
	for i, tt := range tests {
		d, exist, err := TestValues.GetDuration(tt.key)
		if !matchError(t, err, tt.err) {
			t.Errorf("%d: match error: got %v, want %v", i, err, tt.err)
			continue
		}
		if err != nil {
			t.Logf("%d: get duration: %v", i, err)
			continue
		}
		if exist != tt.exist {
			t.Errorf("%d: exist: got %v, want %v", i, exist, tt.exist)
			continue
		}
		if exist {
			if got, want := d, tt.d; got != want {
				t.Errorf("%d: duration: got %v, want %v", i, got, want)
				continue
			}
		}
	}
}

func TestValuesGetSize(t *testing.T) {
	tests := []struct {
		key   string
		err   string
		exist bool
		size  int
	}{
		{key: "nokey", exist: false},
		{key: "k5", exist: true, size: 500 * 1024 * 1024},
		{key: "k6", err: "unknown unit"},
	}
	for i, tt := range tests {
		size, exist, err := TestValues.GetSize(tt.key)
		if !matchError(t, err, tt.err) {
			t.Errorf("%d: match error: got %v, want %v", i, err, tt.err)
			continue
		}
		if err != nil {
			t.Logf("%d: get size: %v", i, err)
			continue
		}
		if exist != tt.exist {
			t.Errorf("%d: exist: got %v, want %v", i, exist, tt.exist)
			continue
		}
		if exist {
			if got, want := size, tt.size; got != want {
				t.Errorf("%d: size: got %v, want %v", i, got, want)
				continue
			}
		}
	}
}
