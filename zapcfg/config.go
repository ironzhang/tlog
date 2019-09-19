package zapcfg

import "go.uber.org/zap/zapcore"

type Level string

type Sink struct {
	Name     string
	MinLevel Level
	MaxLevel Level
	URLs     []string
}

type Encoder struct {
	Name          string
	Encoding      string                `json:"encoding" yaml:"encoding"`
	EncoderConfig zapcore.EncoderConfig `json:"encoderConfig" yaml:"encoderConfig"`
}

type Logger struct {
	Name    string
	Encoder string
	Sinks   []string
}

type Config struct {
	Level    Level
	Sinks    []Sink
	Encoders []Encoder
	Default  Logger
	Loggers  []Logger
}

var Example = Config{
	Level: "InfoLevel",
	Sinks: []Sink{
		{
			Name:     "debug-sink",
			MinLevel: "DebugLevel",
			MaxLevel: "DebugLevel",
			URLs:     []string{"file://log/{{.Process}}/debug.log"},
		},
		{
			Name:     "info-sink",
			MinLevel: "InfoLevel",
			MaxLevel: "FatalLevel",
			URLs:     []string{"file://log/{{.Process}}/info.log"},
		},
		{
			Name:     "warn-sink",
			MinLevel: "WarnLevel",
			MaxLevel: "WarnLevel",
			URLs:     []string{"file://log/{{.Process}}/warn.log"},
		},
		{
			Name:     "error-sink",
			MinLevel: "ErrorLevel",
			MaxLevel: "ErrorLevel",
			URLs:     []string{"file://log/{{.Process}}/error.log"},
		},
		{
			Name:     "fatal-sink",
			MinLevel: "PanicLevel",
			MaxLevel: "FatalLevel",
			URLs:     []string{"file://log/{{.Process}}/fatal.log"},
		},
		{
			Name:     "dirpc-sink",
			MinLevel: "DebugLevel",
			MaxLevel: "FatalLevel",
			URLs:     []string{"file://log/{{.Process}}/dirpc.log"},
		},
		{
			Name:     "access-sink",
			MinLevel: "DebugLevel",
			MaxLevel: "FatalLevel",
			URLs:     []string{"file://log/{{.Process}}/access.log?maxsize=20G"},
		},
		{
			Name:     "event-sink",
			MinLevel: "DebugLevel",
			MaxLevel: "FatalLevel",
			URLs: []string{
				"file://log/{{.Process}}/event.log",
				"kafka://topic.{{.Process}}",
			},
		},
		{
			Name:     "metrics-sink",
			MinLevel: "DebugLevel",
			MaxLevel: "FatalLevel",
			URLs: []string{
				"file://log/{{.Process}}/metrics.log",
				"monitor://172.0.0.1:10000?transport=udp",
			},
		},
	},
	Encoders: []Encoder{
		{
			Name:          "default-encoder",
			Encoding:      "",
			EncoderConfig: zapcore.EncoderConfig{},
		},
	},
	Default: Logger{
		Name:    "",
		Encoder: "default-encoder",
		Sinks:   []string{"debug-sink", "info-sink", "warn-sink", "error-sink", "fatal-sink"},
	},
	Loggers: []Logger{
		{
			Name:    "access",
			Encoder: "default-encoder",
			Sinks:   []string{"access-sink"},
		},
		{
			Name:    "git.xiaojukeji.com/lego/dirpc",
			Encoder: "default-encoder",
			Sinks:   []string{"dirpc-sink"},
		},
	},
}
