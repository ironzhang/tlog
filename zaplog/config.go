package zaplog

import (
	"fmt"
	"sort"

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
	ErrorOutputPaths  []string               `json:"errorOutputPaths" yaml:"errorOutputPaths"`
	Outputs           []Output               `json:"outputs" yaml:"outputs"`
	Categories        []Category             `json:"categories" yaml:"categories"`
	InitialFields     map[string]interface{} `json:"initialFields" yaml:"initialFields"`
}

func (p Output) build() (zapcore.WriteSyncer, func(), error) {
	return zap.Open(p.Path)
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

func (p Config) buildErrorOutputs() (zapcore.WriteSyncer, func(), error) {
	return zap.Open(p.ErrorOutputPaths...)
}

func (p Config) buildOutputs() (map[string]zapcore.WriteSyncer, zapcore.WriteSyncer, func(), error) {
	writers := make(map[string]zapcore.WriteSyncer, len(p.Outputs))
	closers := make([]func(), 0, len(p.Outputs))
	close := func() {
		for _, f := range closers {
			f()
		}
	}

	// 构建输出对象
	for _, out := range p.Outputs {
		sink, closer, err := out.build()
		if err != nil {
			close()
			return nil, nil, nil, err
		}
		closers = append(closers, closer)

		if _, ok := writers[out.Name]; ok {
			close()
			return nil, nil, nil, fmt.Errorf("%q output object is already existed", out.Name)
		}
		writers[out.Name] = sink
	}

	// 构建错误输出对象
	errSink, closer, err := p.buildErrorOutputs()
	if err != nil {
		close()
		return nil, nil, nil, err
	}
	closers = append(closers, closer)

	return writers, errSink, close, nil
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

func (p Config) buildOptions(errOutput zapcore.WriteSyncer) []zap.Option {
	opts := []zap.Option{zap.ErrorOutput(errOutput)}

	if p.Development {
		opts = append(opts, zap.Development())
	}
	if !p.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}

	stackLevel := zap.ErrorLevel
	if p.Development {
		stackLevel = zap.WarnLevel
	}
	if !p.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}

	if len(p.InitialFields) > 0 {
		fs := make([]zap.Field, 0, len(p.InitialFields))
		keys := make([]string, 0, len(p.InitialFields))
		for k := range p.InitialFields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fs = append(fs, zap.Any(k, p.InitialFields[k]))
		}
		opts = append(opts, zap.Fields(fs...))
	}

	return opts
}

func (p Config) build(opts ...zap.Option) (*zap.Logger, func(), error) {
	outputs, errOutput, close, err := p.buildOutputs()
	if err != nil {
		return nil, nil, err
	}
	core, err := p.buildCore(outputs)
	if err != nil {
		close()
		return nil, nil, err
	}
	log := zap.New(core, p.buildOptions(errOutput)...)
	if len(opts) > 0 {
		log = log.WithOptions(opts...)
	}
	return log, close, nil
}

func (p Config) BuildZapLogger(options ...zap.Option) (*zap.Logger, error) {
	logger, _, err := p.build(options...)
	return logger, err
}

func (p Config) Build() (*Logger, error) {
	return nil, nil
}
