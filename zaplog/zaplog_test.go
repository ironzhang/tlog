package zaplog

import (
	"testing"

	"github.com/ironzhang/tlog/iface"
)

func TestLoggerLevel(t *testing.T) {
	logger, err := New(NewDevelopmentConfig())
	if err != nil {
		t.Fatalf("new: %v", err)
	}

	if got, want := logger.GetLevel(), iface.DEBUG; got != want {
		t.Errorf("level: got %v, want %v", got, want)
	}
	logger.SetLevel(iface.INFO)
	if got, want := logger.GetLevel(), iface.INFO; got != want {
		t.Errorf("level: got %v, want %v", got, want)
	}
}

func TestLoggerSync(t *testing.T) {
	tsink := RegisterTestSink(t, "TestLoggerSync")
	cfg := Config{
		Level: iface.DEBUG,
		Cores: []CoreConfig{
			{
				Name:     "Test",
				MinLevel: iface.DEBUG,
				MaxLevel: iface.FATAL,
				URLs:     []string{"TestLoggerSync://1", "TestLoggerSync://2"},
			},
		},
		Loggers: []LoggerConfig{
			{
				Cores: []string{"Test"},
			},
		},
	}

	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	if err := logger.Sync(); err != nil {
		t.Fatalf("sync: %v", err)
	}
	if got, want := tsink.syncCount, 2; got != want {
		t.Errorf("sync count: got %v, want %v", got, want)
	}
}

func TestLoggerClose(t *testing.T) {
	tsink := RegisterTestSink(t, "TestLoggerClose")
	cfg := Config{
		Level: iface.DEBUG,
		Cores: []CoreConfig{
			{
				Name:     "Test",
				MinLevel: iface.DEBUG,
				MaxLevel: iface.FATAL,
				URLs:     []string{"TestLoggerClose://1", "TestLoggerClose://2"},
			},
		},
		Loggers: []LoggerConfig{
			{
				Cores: []string{"Test"},
			},
		},
	}

	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	if err := logger.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	if got, want := tsink.closeCount, 2; got != want {
		t.Errorf("close count: got %v, want %v", got, want)
	}
}

//func TestNamed(t *testing.T) {
//	cfg := Config{
//		Loggers: []LoggerConfig{
//			{Name: "root"},
//			{Name: "n1"},
//			{Name: "n2"},
//		},
//	}
//	logger, err := New(cfg)
//	if err != nil {
//		t.Fatalf("new: %v", err)
//	}
//
//	tests := []struct {
//		name   string
//		expect string
//	}{
//		{name: "root", expect: "root"},
//		{name: "n1", expect: "n1"},
//		{name: "n2", expect: "n2"},
//		//{name: "n3", expect: "root"},
//	}
//	for i, tt := range tests {
//		if got, want := logger.Named(tt.name), logger.Named(tt.expect); got != want {
//			t.Errorf("%d: get logger: got %p, want %p", i, got, want)
//		}
//	}
//}
