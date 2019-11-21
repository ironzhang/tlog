package main

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/ironzhang/tlog"
	"github.com/ironzhang/tlog/iface"
	"github.com/ironzhang/tlog/zaplog"
)

var configs = []zaplog.Config{
	{
		Level: iface.DEBUG,
		Sinks: []zaplog.SinkConfig{
			{
				Name: "StderrSink",
				URL:  "stderr",
			},
		},
		Cores: []zaplog.CoreConfig{
			{
				Name:     "StderrCore",
				Encoding: "console",
				Encoder: zaplog.EncoderConfig{
					MessageKey:     "M",
					LevelKey:       "L",
					TimeKey:        "T",
					NameKey:        "N",
					CallerKey:      "C",
					StacktraceKey:  "S",
					EncodeLevel:    zaplog.CapitalLevelEncoder,
					EncodeTime:     zaplog.ISO8601TimeEncoder,
					EncodeDuration: zaplog.StringDurationEncoder,
					EncodeCaller:   zaplog.ShortCallerEncoder,
				},
				MinLevel: iface.DEBUG,
				MaxLevel: iface.FATAL,
				Sinks:    []string{"StderrSink"},
			},
		},
		Loggers: []zaplog.LoggerConfig{
			{
				Name:            "",
				DisableCaller:   false,
				StacktraceLevel: zaplog.DisableStacktrace,
				Cores:           []string{"StderrCore"},
			},
		},
	},
	{
		Level: iface.INFO,
		Sinks: []zaplog.SinkConfig{
			{
				Name: "debugSink",
				URL:  "./log/debug.log",
			},
			{
				Name: "infoSink",
				URL:  "./log/info.log",
			},
			{
				Name: "warnSink",
				URL:  "./log/warn.log",
			},
			{
				Name: "errorSink",
				URL:  "./log/error.log",
			},
			{
				Name: "fatalSink",
				URL:  "./log/fatal.log",
			},
			{
				Name: "accessSink",
				URL:  "./log/access.log",
			},
		},
		Cores: []zaplog.CoreConfig{
			{
				Name:     "debugCore",
				Encoding: "console",
				Encoder:  zaplog.EncoderConfig{},
				MinLevel: iface.DEBUG,
				MaxLevel: iface.DEBUG,
				Sinks:    []string{"debugSink"},
			},
			{
				Name:     "infoCore",
				Encoding: "console",
				Encoder:  zaplog.EncoderConfig{},
				MinLevel: iface.INFO,
				MaxLevel: iface.FATAL,
				Sinks:    []string{"infoSink"},
			},
			{
				Name:     "warnCore",
				Encoding: "console",
				Encoder:  zaplog.EncoderConfig{},
				MinLevel: iface.WARN,
				MaxLevel: iface.FATAL,
				Sinks:    []string{"warnSink"},
			},
			{
				Name:     "errorCore",
				Encoding: "console",
				Encoder:  zaplog.EncoderConfig{},
				MinLevel: iface.ERROR,
				MaxLevel: iface.FATAL,
				Sinks:    []string{"errorSink"},
			},
			{
				Name:     "fatalCore",
				Encoding: "console",
				Encoder:  zaplog.EncoderConfig{},
				MinLevel: iface.PANIC,
				MaxLevel: iface.FATAL,
				Sinks:    []string{"fatalSink"},
			},
			{
				Name:     "accessCore",
				Encoding: "console",
				Encoder:  zaplog.EncoderConfig{},
				MinLevel: iface.DEBUG,
				MaxLevel: iface.FATAL,
				Sinks:    []string{"accessSink"},
			},
		},
		Loggers: []zaplog.LoggerConfig{
			{
				Name:            "",
				DisableCaller:   false,
				StacktraceLevel: zaplog.DisableStacktrace,
				Cores:           []string{"debugCore", "infoCore", "warnCore", "errorCore", "fatalCore"},
			},
			{
				Name:            "access",
				DisableCaller:   false,
				StacktraceLevel: zaplog.DisableStacktrace,
				Cores:           []string{"accessCore"},
			},
		},
	},
}

type MarshalFunc func(v interface{}) ([]byte, error)

func WriteConfig(marshal MarshalFunc, file string, cfg zaplog.Config) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := marshal(cfg)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func MarshalJSON(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "\t")
}

func main() {
	clean := false
	if len(os.Args) >= 2 && os.Args[1] == "clean" {
		clean = true
	}

	entries := []struct {
		file    string
		marshal MarshalFunc
		config  zaplog.Config
	}{
		{
			file:    "./std.json",
			marshal: MarshalJSON,
			config:  configs[0],
		},
		{
			file:    "./std.yaml",
			marshal: yaml.Marshal,
			config:  configs[0],
		},
		{
			file:    "./example.json",
			marshal: MarshalJSON,
			config:  configs[1],
		},
		{
			file:    "./example.yaml",
			marshal: yaml.Marshal,
			config:  configs[1],
		},
	}

	if clean {
		for _, e := range entries {
			os.Remove(e.file)
		}
	} else {
		for _, e := range entries {
			if err := WriteConfig(e.marshal, e.file, e.config); err != nil {
				tlog.Errorf("write %q config: %v", e.file, err)
			}
		}
	}
}
