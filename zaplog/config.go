package zaplog

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type StreamConfig struct {
	Name     string   `json:"name" yaml:"name"`
	MinLevel Level    `json:"minLevel" yaml:"minLevel"`
	MaxLevel Level    `json:"maxLevel" yaml:"maxLevel"`
	URLs     []string `json:"urls" yaml:"urls"`
}

type EncoderConfig struct {
	Name          string `json:"name" yaml:"name"`
	Encoding      string `json:"encoding" yaml:"encoding"`
	MessageKey    string `json:"messageKey" yaml:"messageKey"`
	LevelKey      string `json:"levelKey" yaml:"levelKey"`
	TimeKey       string `json:"timeKey" yaml:"timeKey"`
	NameKey       string `json:"nameKey" yaml:"nameKey"`
	CallerKey     string `json:"callerKey" yaml:"callerKey"`
	StacktraceKey string `json:"stacktraceKey" yaml:"stacktraceKey"`
	LineEnding    string `json:"lineEnding" yaml:"lineEnding"`
	//	EncodeLevel    zapcore.LevelEncoder    `json:"levelEncoder" yaml:"levelEncoder"`
	//	EncodeTime     zapcore.TimeEncoder     `json:"timeEncoder" yaml:"timeEncoder"`
	//	EncodeDuration zapcore.DurationEncoder `json:"durationEncoder" yaml:"durationEncoder"`
	//	EncodeCaller   zapcore.CallerEncoder   `json:"callerEncoder" yaml:"callerEncoder"`
	//	EncodeName     zapcore.NameEncoder     `json:"nameEncoder" yaml:"nameEncoder"`
}

type LoggerConfig struct {
	Name    string   `json:"name" yaml:"name"`
	Encoder string   `json:"encoder" yaml:"encoder"`
	Streams []string `json:"streams" yaml:"streams"`
}

type Config struct {
	Level    Level           `json:"level" yaml:"level"`
	Streams  []StreamConfig  `json:"streams" yaml:"streams"`
	Encoders []EncoderConfig `json:"encoders" yaml:"encoders"`
	Loggers  []LoggerConfig  `json:"loggers" yaml:"loggers"`
}

func LoadConfig(file string) (Config, error) {
	readFile := readJSON
	switch strings.ToUpper(filepath.Ext(file)) {
	case ".JSON":
		readFile = readJSON
	case ".YAML":
		readFile = readYAML
	}

	var cfg Config
	if err := readFile(file, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func WriteConfig(file string, cfg Config) error {
	switch strings.ToUpper(filepath.Ext(file)) {
	case ".JSON":
		return writeJSON(file, cfg)
	case ".YAML":
		return writeYAML(file, cfg)
	}
	return writeJSON(file, cfg)
}

func readJSON(file string, a interface{}) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, a)
}

func writeJSON(file string, a interface{}) error {
	data, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, data, 0666)
}

func readYAML(file string, a interface{}) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, a)
}

func writeYAML(file string, a interface{}) error {
	data, err := yaml.Marshal(a)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, data, 0666)
}
