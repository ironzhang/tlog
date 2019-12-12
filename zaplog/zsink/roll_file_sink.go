package zsink

import (
	"fmt"
	"net/url"
	"strings"

	"go.uber.org/zap"

	"github.com/ironzhang/tlog/zaplog/zsink/rollfile"
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

	suffix, ok := params.Get("suffix")
	if ok {
		layout, err := suffixToLayout(suffix)
		if err != nil {
			return nil, err
		}
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

	maxSize, ok, err := params.GetSize("maxSize")
	if err != nil {
		return nil, err
	}
	if ok {
		opts = append(opts, rollfile.SetMaxSize(maxSize))
	}

	return opts, nil
}

func suffixToLayout(s string) (string, error) {
	switch strings.ToLower(s) {
	case "d", "day":
		return rollfile.DayLayout, nil
	case "h", "hour":
		return rollfile.HourLayout, nil
	case "s", "second":
		return rollfile.SecondLayout, nil
	case "n", "nano":
		return rollfile.NanoLayout, nil
	default:
		return "", fmt.Errorf("unknown suffix pattern %q", s)
	}
}

func init() {
	if err := zap.RegisterSink("rfile", newRollFileSink); err != nil {
		panic(err)
	}
}
