package zaplog

import (
	"os"
	"reflect"
	"testing"
)

var ComplexConfig = Config{
	Level: INFO,
	Streams: []StreamConfig{
		{
			Name:     "Trace-Stream",
			MinLevel: TRACE,
			MaxLevel: TRACE,
			URLs:     []string{"file://log/{{.Process}}/debug.log"},
		},
		{
			Name:     "Debug-Stream",
			MinLevel: DEBUG,
			MaxLevel: DEBUG,
			URLs:     []string{"file://log/{{.Process}}/debug.log"},
		},
		{
			Name:     "Info-Stream",
			MinLevel: INFO,
			MaxLevel: FATAL,
			URLs:     []string{"file://log/{{.Process}}/info.log"},
		},
		{
			Name:     "Warn-Stream",
			MinLevel: WARN,
			MaxLevel: WARN,
			URLs:     []string{"file://log/{{.Process}}/warn.log"},
		},
		{
			Name:     "Error-Stream",
			MinLevel: ERROR,
			MaxLevel: ERROR,
			URLs:     []string{"file://log/{{.Process}}/error.log"},
		},
		{
			Name:     "Fatal-Stream",
			MinLevel: PANIC,
			MaxLevel: FATAL,
			URLs:     []string{"file://log/{{.Process}}/fatal.log"},
		},
		{
			Name: "Zerone-Stream",
			URLs: []string{"file://log/{{.Process}}/zerone.log"},
		},
		{
			Name: "Access-Stream",
			URLs: []string{"file://log/{{.Process}}/access.log?maxsize=20G"},
		},
		{
			Name: "Event-Stream",
			URLs: []string{
				"file://log/{{.Process}}/event.log",
				"kafka://topic.{{.Process}}",
			},
		},
	},
	Encoders: []EncoderConfig{
		{
			Name: "Access-Encoder",
		},
	},
	Loggers: []LoggerConfig{
		{
			Name:    "Default",
			Encoder: "",
			Streams: []string{"Trace-Stream", "Debug-Stream", "Info-Stream", "Warn-Stream", "Error-Stream", "Fatal-Stream"},
		},
		{
			Name:    "access",
			Encoder: "Access-Encoder",
			Streams: []string{"access-sink"},
		},
		{
			Name:    "github.com/ironzhang/zerone",
			Encoder: "",
			Streams: []string{"Zerone-Stream"},
		},
	},
}

func TestWriteConfig(t *testing.T) {
	tests := []struct {
		file string
		conf Config
	}{
		{
			file: "defalut.config.yaml",
			conf: Config{
				Streams:  []StreamConfig{},
				Encoders: []EncoderConfig{},
				Loggers:  []LoggerConfig{},
			},
		},
		{file: "example.config.yaml", conf: ExampleConfig},
		{file: "complex.config.json", conf: ComplexConfig},
	}
	for i, tt := range tests {
		if err := WriteConfig(tt.file, tt.conf); err != nil {
			t.Errorf("%d: write config: %v", i, err)
			continue
		}
		cfg, err := LoadConfig(tt.file)
		if err != nil {
			t.Errorf("%d: load config: %v", i, err)
			continue
		}
		if got, want := cfg, tt.conf; !reflect.DeepEqual(got, want) {
			t.Errorf("%d: got %v, want %v", i, got, want)
			continue
		}
		os.Remove(tt.file)
	}
}
