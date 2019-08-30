package zapx

import (
	"fmt"
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Output struct {
	Name string `json:"name" yaml:"name"`
	Path string `json:"path" yaml:"path"`
}

type Category struct {
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
	Categories        []Category             `json:"categories" yaml:"categories"`
	ErrorOutputPaths  []string               `json:"errorOutputPaths" yaml:"errorOutputPaths"`
	InitialFields     map[string]interface{} `json:"initialFields" yaml:"initialFields"`
}

func (p Output) build() (zap.Sink, error) {
	return newSink(p.Name, p.Path)
}

func (p Category) build(lvl zap.AtomicLevel, outputs map[string]zapcore.WriteSyncer) (zapcore.Core, error) {
	// 构建日志编码器
	enc, err := newEncoder(p.Encoding, p.EncoderConfig)
	if err != nil {
		return nil, err
	}

	// 构建日志输出对象列表
	ws := make([]zapcore.WriteSyncer, 0, len(p.Outputs))
	for _, name := range p.Outputs {
		w, ok := outputs[name]
		if !ok {
			return nil, fmt.Errorf("%q output object is not existed", name)
		}
		ws = append(ws, w)
	}

	// 构建日志分级输出规则
	enab := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		if l < p.MinLevel || l > p.MaxLevel {
			return false
		}
		return lvl.Enabled(l)
	})

	return zapcore.NewCore(enc, zap.CombineWriteSyncers(ws...), enab), nil
}

func (p Config) buildOutputs() (map[string]zapcore.WriteSyncer, func(), error) {
	writers := make(map[string]zapcore.WriteSyncer, len(p.Outputs))
	closers := make([]io.Closer, 0, len(p.Outputs))
	close := func() {
		for _, c := range closers {
			c.Close()
		}
	}

	for _, out := range p.Outputs {
		sink, err := out.build()
		if err != nil {
			close()
			return nil, nil, err
		}
		closers = append(closers, sink)

		if _, ok := writers[out.Name]; ok {
			close()
			return nil, nil, fmt.Errorf("%q output object is already existed", out.Name)
		}
		writers[out.Name] = sink
	}

	return writers, close, nil
}

func (p Config) buildCore(outputs map[string]zapcore.WriteSyncer) (zapcore.Core, error) {
	cores := make([]zapcore.Core, 0, len(p.Categories))
	for _, cate := range p.Categories {
		core, err := cate.build(p.Level, outputs)
		if err != nil {
			return nil, err
		}
		cores = append(cores, core)
	}
	return zapcore.NewTee(cores...), nil
}

func (p Config) build(options ...zap.Option) (*zap.Logger, func(), error) {
	outputs, close, err := p.buildOutputs()
	if err != nil {
		return nil, nil, err
	}
	core, err := p.buildCore(outputs)
	if err != nil {
		close()
		return nil, nil, err
	}
	return zap.New(core, options...), close, nil
}
