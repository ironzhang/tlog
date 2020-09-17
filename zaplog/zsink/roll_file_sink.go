package zsink

import (
	"fmt"
	"net/url"
	"strings"

	"go.uber.org/zap"

	"git.xiaojukeji.com/pearls/tlog/zaplog/zsink/rollfile"
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
	case "localhost", "rootdir", "$localhost", "$rootdir":
		return u.Path, nil
	case "workdir", "$workdir":
		return strings.TrimPrefix(u.Path, "/"), nil
	default:
		return "", fmt.Errorf("invalid hostname %q", hostname)
	}
}

func parseFileOptions(u *url.URL) (opts []rollfile.Option, err error) {
	params := values(u.Query())

	cut, ok := params.Get("cut")
	if ok {
		cutfmt, err := stringToCutFormat(cut)
		if err != nil {
			return nil, err
		}
		opts = append(opts, rollfile.SetCutFormat(cutfmt))
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

func stringToCutFormat(s string) (rollfile.CutFormat, error) {
	switch strings.ToLower(s) {
	case "size":
		return rollfile.SizeCut, nil
	case "h", "hour":
		return rollfile.HourCut, nil
	case "d", "day":
		return rollfile.DayCut, nil
	default:
		return "", fmt.Errorf("unknown cut format %q", s)
	}
}

func init() {
	if err := zap.RegisterSink("rfile", newRollFileSink); err != nil {
		panic(err)
	}
}
