package iface

import "testing"

func TestLevelString(t *testing.T) {
	tests := []struct {
		lvl Level
		str string
	}{
		{lvl: DEBUG, str: "debug"},
		{lvl: INFO, str: "info"},
		{lvl: WARN, str: "warn"},
		{lvl: ERROR, str: "error"},
		{lvl: PANIC, str: "panic"},
		{lvl: FATAL, str: "fatal"},
		{lvl: -10, str: "Level(-10)"},
	}
	for i, tt := range tests {
		if got, want := tt.lvl.String(), tt.str; got != want {
			t.Errorf("%d: level string: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: level string: got %v", i, tt.lvl.String())
	}
}

func TestLevelMarshalText(t *testing.T) {
	tests := []struct {
		level Level
		text  string
	}{
		{level: DEBUG, text: "debug"},
		{level: INFO, text: "info"},
		{level: WARN, text: "warn"},
		{level: ERROR, text: "error"},
		{level: PANIC, text: "panic"},
		{level: FATAL, text: "fatal"},
	}
	for i, tt := range tests {
		text, err := tt.level.MarshalText()
		if err != nil {
			t.Errorf("%d: level marshal text: %v", i, err)
			continue
		}
		if got, want := string(text), tt.text; got != want {
			t.Errorf("%d: level text: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: level text: got %s", i, text)
	}
}

func TestLevelUnmarshalText(t *testing.T) {
	tests := []struct {
		text  string
		level Level
	}{
		{text: "debug", level: DEBUG},
		{text: "info", level: INFO},
		{text: "warn", level: WARN},
		{text: "error", level: ERROR},
		{text: "panic", level: PANIC},
		{text: "fatal", level: FATAL},
		{text: "Debug", level: DEBUG},
		{text: "DEBUG", level: DEBUG},
		{text: "", level: INFO},
	}
	for i, tt := range tests {
		var level Level
		err := level.UnmarshalText([]byte(tt.text))
		if err != nil {
			t.Errorf("%d: level unmarshal text: %v", i, err)
			continue
		}
		if got, want := level, tt.level; got != want {
			t.Errorf("%d: level: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: level: got %v", i, level)
	}
}

func TestStringToLevel(t *testing.T) {
	tests := []struct {
		text  string
		level Level
	}{
		{text: "debug", level: DEBUG},
		{text: "info", level: INFO},
		{text: "warn", level: WARN},
		{text: "error", level: ERROR},
		{text: "panic", level: PANIC},
		{text: "fatal", level: FATAL},
		{text: "Debug", level: DEBUG},
		{text: "DEBUG", level: DEBUG},
		{text: "", level: INFO},
	}
	for i, tt := range tests {
		level, err := StringToLevel(tt.text)
		if err != nil {
			t.Errorf("%d: string to level: %v", i, err)
			continue
		}
		if got, want := level, tt.level; got != want {
			t.Errorf("%d: level: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: level: got %v", i, level)
	}
}
