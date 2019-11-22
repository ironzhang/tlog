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
		Cores: []zaplog.CoreConfig{
			{
				Name:     "Stderr",
				Encoding: "console",
				Encoder: zaplog.EncoderConfig{
					MessageKey:     "M",
					LevelKey:       "L",
					TimeKey:        "T",
					NameKey:        "N",
					CallerKey:      "C",
					StacktraceKey:  "S",
					EncodeLevel:    zaplog.CapitalColorLevelEncoder,
					EncodeTime:     zaplog.ISO8601TimeEncoder,
					EncodeDuration: zaplog.StringDurationEncoder,
					EncodeCaller:   zaplog.ShortCallerEncoder,
					EncodeName:     zaplog.FullNameEncoder,
				},
				MinLevel: iface.DEBUG,
				MaxLevel: iface.FATAL,
				URLs:     []string{"stderr"},
			},
		},
		Loggers: []zaplog.LoggerConfig{
			{
				Name:            "",
				DisableCaller:   false,
				StacktraceLevel: zaplog.DisableStacktrace,
				Cores:           []string{"Stderr"},
			},
		},
	},
	{
		Level: iface.INFO,
		Cores: []zaplog.CoreConfig{
			{
				Name:     "Debug",
				Encoding: "",
				Encoder:  zaplog.EncoderConfig{},
				MinLevel: iface.DEBUG,
				MaxLevel: iface.DEBUG,
				URLs:     []string{"./log/debug.log"},
			},
			{
				Name:     "Info",
				Encoding: "",
				Encoder:  zaplog.EncoderConfig{},
				MinLevel: iface.INFO,
				MaxLevel: iface.FATAL,
				URLs:     []string{"./log/info.log"},
			},
			{
				Name:     "Warn",
				Encoding: "",
				Encoder:  zaplog.EncoderConfig{},
				MinLevel: iface.WARN,
				MaxLevel: iface.FATAL,
				URLs:     []string{"./log/warn.log"},
			},
			{
				Name:     "Error",
				Encoding: "",
				Encoder:  zaplog.EncoderConfig{},
				MinLevel: iface.ERROR,
				MaxLevel: iface.FATAL,
				URLs:     []string{"./log/error.log"},
			},
			{
				Name:     "Fatal",
				Encoding: "",
				Encoder:  zaplog.EncoderConfig{},
				MinLevel: iface.PANIC,
				MaxLevel: iface.FATAL,
				URLs:     []string{"./log/fatal.log"},
			},
			{
				Name:     "Access",
				Encoding: "",
				Encoder:  zaplog.EncoderConfig{},
				MinLevel: iface.DEBUG,
				MaxLevel: iface.FATAL,
				URLs:     []string{"./log/access.log"},
			},
		},
		Loggers: []zaplog.LoggerConfig{
			{
				Name:            "",
				DisableCaller:   false,
				StacktraceLevel: zaplog.DisableStacktrace,
				Cores:           []string{"Debug", "Info", "Warn", "Error", "Fatal"},
			},
			{
				Name:            "access",
				DisableCaller:   false,
				StacktraceLevel: zaplog.DisableStacktrace,
				Cores:           []string{"Access"},
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
