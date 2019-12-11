package zsink

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ironzhang/tlog/zaplog/zsink/rollfile"
	"go.uber.org/zap"
)

func newRollFileSink(u *url.URL) (zap.Sink, error) {
	filename, err := parseFilePath(u)
	if err != nil {
		return nil, err
	}
	opts, err := parseFileOptions(u)
	if err != nil {
		return nil, err
	}
	file, err := rollfile.Open(filename, opts...)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func parseFilePath(u *url.URL) (string, error) {
	switch hostname := u.Hostname(); hostname {
	case "localhost", "rootdir":
		return u.Path, nil
	case "workdir":
		return strings.TrimPrefix(u.Path, "/"), nil
	default:
		return "", fmt.Errorf("invalid hostname %q", hostname)
	}
}

func parseFileOptions(u *url.URL) (opts []rollfile.Option, err error) {
	params := values(u.Query())

	layout, ok := params.Get("layout")
	if ok {
		opts = append(opts, rollfile.SetLayout(layout))
	}

	period, ok, err := params.GetDuration("period")
	if err != nil {
		return nil, err
	}
	if ok {
		opts = append(opts, rollfile.SetPeriod(period))
	}

	maxSeq, ok, err := params.GetInt("maxSeq")
	if err != nil {
		return nil, err
	}
	if ok {
		opts = append(opts, rollfile.SetMaxSeq(maxSeq))
	}

	maxSize, ok, err := params.GetInt("maxSize")
	if err != nil {
		return nil, err
	}
	if ok {
		opts = append(opts, rollfile.SetMaxSize(maxSize))
	}

	return opts, nil
}

type values url.Values

func (v values) Get(key string) (string, bool) {
	if v == nil {
		return "", false
	}
	vs := v[key]
	if len(vs) == 0 {
		return "", false
	}
	return vs[0], true
}

func (v values) GetInt(key string) (int, bool, error) {
	s, ok := v.Get(key)
	if !ok {
		return 0, false, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, false, err
	}
	return n, true, nil
}

func (v values) GetDuration(key string) (time.Duration, bool, error) {
	s, ok := v.Get(key)
	if !ok {
		return 0, false, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, false, err
	}
	return d, true, nil
}

func init() {
	if err := zap.RegisterSink("rfile", newRollFileSink); err != nil {
		panic(err)
	}
}
