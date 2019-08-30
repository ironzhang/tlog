package zapx

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Output struct {
	Name string `json:"name" yaml:"name"`
	Path string `json:"path" yaml:"path"`
}

type Catalog struct {
	MinLevel      zapcore.Level         `json:"minLevel" yaml:"minLevel"`
	MaxLevel      zapcore.Level         `json:"maxLevel" yaml:"maxLevel"`
	Encoding      string                `json:"encoding" yaml:"encoding"`
	EncoderConfig zapcore.EncoderConfig `json:"encoderConfig" yaml:"encoderConfig"`
	Outputs       []string
}

type Config struct {
	Level             zap.AtomicLevel        `json:"level" yaml:"level"`
	Development       bool                   `json:"development" yaml:"development"`
	DisableCaller     bool                   `json:"disableCaller" yaml:"disableCaller"`
	DisableStacktrace bool                   `json:"disableStacktrace" yaml:"disableStacktrace"`
	Outputs           []Output               `json:"outputs" yaml:"outputs"`
	Catalogs          []Catalog              `json:"catalog" yaml:"catalog"`
	ErrorOutputPaths  []string               `json:"errorOutputPaths" yaml:"errorOutputPaths"`
	InitialFields     map[string]interface{} `json:"initialFields" yaml:"initialFields"`
}

func (p Output) build() (Sink, error) {
	return newSink(p.Name, p.Path)
}

func (p Catalog) build(lvl zap.AtomicLevel, outputs map[string]zapcore.WriteSyncer) (zapcore.Core, error) {
	ws := make([]zapcore.WriteSyncer, 0, len(p.Outputs))
	for _, name := range p.Outputs {
		w, ok := outputs[name]
		if !ok {
			return nil, fmt.Errorf("%q output object is not existed", name)
		}
		ws = append(ws, w)
	}

	enc, err := newEncoder(p.Encoding, p.EncoderConfig)
	if err != nil {
		return nil, err
	}

	enab := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		if l < p.MinLevel || l > p.MaxLevel {
			return false
		}
		return lvl.Enabled(l)
	})

	return zapcore.NewCore(enc, zap.CombineWriteSyncers(ws...), enab), nil
}

func (p Config) buildOutputs() (map[string]Sink, error) {
	sinks := make(map[string]Sink, len(p.Outputs))
	close := func() {
		for _, s := range sinks {
			s.Close()
		}
	}

	for _, out := range p.Outputs {
		sink, err := out.build()
		if err != nil {
			close()
			return nil, err
		}
		if _, ok := sinks[out.Name]; ok {
			sink.Close()
			close()
			return nil, fmt.Errorf("%q output object is already existed", out.Name)
		}
		sinks[out.Name] = sink
	}

	return sinks, nil
}

func (p Config) build() (*zap.Logger, error) {
	return nil, nil
}
