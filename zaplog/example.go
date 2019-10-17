package zaplog

var ExampleConfig = Config{
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
	},
	Encoders: []EncoderConfig{},
	Loggers: []LoggerConfig{
		{
			Name:    "Default",
			Encoder: "",
			Streams: []string{"Trace-Stream", "Debug-Stream", "Info-Stream", "Warn-Stream", "Error-Stream", "Fatal-Stream"},
		},
	},
}
