package zsink

import (
	"net/url"
	"strconv"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
)

type rollFileSink struct {
	lumberjack.Logger
}

func (p *rollFileSink) Sync() error {
	return nil
}

func newRollFileSink(u *url.URL) (zap.Sink, error) {
	values := u.Query()
	maxSize, err := getIntValue(values, "maxSize")
	if err != nil {
		return nil, err
	}
	maxAge, err := getIntValue(values, "maxAge")
	if err != nil {
		return nil, err
	}

	println(maxSize)

	return &rollFileSink{
		Logger: lumberjack.Logger{
			Filename: u.Path[1:],
			MaxSize:  maxSize,
			MaxAge:   maxAge,
		},
	}, nil
}

func getIntValue(values url.Values, name string) (int, error) {
	s := values.Get(name)
	if s == "" {
		return 0, nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func init() {
	zap.RegisterSink("rollfile", newRollFileSink)
}
