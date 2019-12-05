package main

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/ironzhang/tlog"
	"github.com/ironzhang/tlog/iface"
	"github.com/ironzhang/tlog/zaplog"
)

func NewAccessLoggerConfig() zaplog.Config {
	return zaplog.Config{
		Level: iface.INFO,
		Loggers: []zaplog.LoggerConfig{
			{
				Name:            "",
				DisableCaller:   false,
				StacktraceLevel: zaplog.PanicStacktrace,
				Encoding:        "",
				Encoder:         zaplog.NewConsoleEncoderConfig(),
				Outputs: []zaplog.OutputConfig{
					{MinLevel: iface.DEBUG, MaxLevel: iface.FATAL, URLs: []string{"./log/debug.log"}},
					{MinLevel: iface.INFO, MaxLevel: iface.FATAL, URLs: []string{"./log/info.log"}},
					{MinLevel: iface.WARN, MaxLevel: iface.FATAL, URLs: []string{"./log/warn.log"}},
					{MinLevel: iface.ERROR, MaxLevel: iface.FATAL, URLs: []string{"./log/error.log"}},
					{MinLevel: iface.PANIC, MaxLevel: iface.FATAL, URLs: []string{"./log/fatal.log"}},
				},
			},
			{
				Name:            "access",
				DisableCaller:   false,
				StacktraceLevel: zaplog.PanicStacktrace,
				Encoding:        "",
				Encoder:         zaplog.NewConsoleEncoderConfig(),
				Outputs: []zaplog.OutputConfig{
					{MinLevel: iface.DEBUG, MaxLevel: iface.FATAL, URLs: []string{"./log/access.log"}},
				},
			},
		},
	}
}

var configs = []zaplog.Config{
	zaplog.NewDevelopmentConfig(),
	zaplog.NewProductionConfig(),
	NewAccessLoggerConfig(),
}

type MarshalFunc func(v interface{}) ([]byte, error)

func MarshalJSON(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "\t")
}

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
			file:    "./development.json",
			marshal: MarshalJSON,
			config:  configs[0],
		},
		{
			file:    "./development.yaml",
			marshal: yaml.Marshal,
			config:  configs[0],
		},
		{
			file:    "./production.json",
			marshal: MarshalJSON,
			config:  configs[1],
		},
		{
			file:    "./production.yaml",
			marshal: yaml.Marshal,
			config:  configs[1],
		},
		{
			file:    "./access.json",
			marshal: MarshalJSON,
			config:  configs[2],
		},
		{
			file:    "./access.yaml",
			marshal: yaml.Marshal,
			config:  configs[2],
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
