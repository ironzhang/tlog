package zsink

import (
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"git.xiaojukeji.com/pearls/tlog/zaplog/zsink/rollfile"
)

func ParseTestURL(t *testing.T, rawurl string) *url.URL {
	u, err := url.Parse(rawurl)
	if err != nil {
		t.Fatalf("parse url: %v", err)
	}
	return u
}

func TestParseFilePath(t *testing.T) {
	tests := []struct {
		url  *url.URL
		err  string
		path string
	}{
		{
			url:  ParseTestURL(t, "rfile://localhost/log/a.log"),
			path: "/log/a.log",
		},
		{
			url:  ParseTestURL(t, "rfile://rootdir/log/a.log"),
			path: "/log/a.log",
		},
		{
			url:  ParseTestURL(t, "rfile://workdir/log/a.log"),
			path: "log/a.log",
		},
		{
			url: ParseTestURL(t, "rfile:///log/a.log"),
			err: "invalid hostname",
		},
	}
	for i, tt := range tests {
		path, err := parseFilePath(tt.url)
		if !matchError(t, err, tt.err) {
			t.Errorf("%d: match error: got %v, want %v", i, err, tt.err)
			continue
		}
		if err != nil {
			t.Logf("%d: parse file path: %v", i, err)
			continue
		}
		if got, want := path, tt.path; got != want {
			t.Errorf("%d: path: got %v, want %v", i, got, want)
			continue
		}
	}
}

func TestStringToCutFormat(t *testing.T) {
	tests := []struct {
		suffix string
		format rollfile.CutFormat
		err    string
	}{
		{suffix: "d", format: rollfile.DayCut},
		{suffix: "day", format: rollfile.DayCut},
		{suffix: "h", format: rollfile.HourCut},
		{suffix: "hour", format: rollfile.HourCut},
		{suffix: "Day", format: rollfile.DayCut},
		{suffix: "DAY", format: rollfile.DayCut},
		{suffix: "days", err: "unknown cut format"},
	}
	for i, tt := range tests {
		format, err := stringToCutFormat(tt.suffix)
		if !matchError(t, err, tt.err) {
			t.Errorf("%d: match error: got %v, want %v", i, err, tt.err)
			continue
		}
		if err != nil {
			t.Logf("%d: string to cut format: %v", i, err)
			continue
		}
		if got, want := format, tt.format; got != want {
			t.Errorf("%d: format: got %v, want %v", i, got, want)
			continue
		}
	}
}

func TestNewNewRollFileSink(t *testing.T) {
	tests := []struct {
		url *url.URL
		err string
	}{
		{url: ParseTestURL(t, "rfile://workdir/testdata/log/a.log")},
		{url: ParseTestURL(t, "rfile://workdir2/testdata/log/a.log"), err: "invalid hostname"},
		{url: ParseTestURL(t, "rfile://workdir/testdata/log/b.log?cut=hour")},
		{url: ParseTestURL(t, "rfile://workdir/testdata/log/f.log?cut=year"), err: "unknown cut format"},
		{url: ParseTestURL(t, "rfile://workdir/testdata/log/d.log?maxSeq=10")},
		{url: ParseTestURL(t, "rfile://workdir/testdata/log/d.log?maxSeq=a"), err: "strconv.Atoi"},
		{url: ParseTestURL(t, "rfile://workdir/testdata/log/f.log?maxSize=1G")},
		{url: ParseTestURL(t, "rfile://workdir/testdata/log/f.log?maxSize=1T"), err: "unknown unit"},
	}
	for i, tt := range tests {
		sink, err := newRollFileSink(tt.url)
		if !matchError(t, err, tt.err) {
			t.Errorf("%d: match error: %v", i, err)
			continue
		}
		if err != nil {
			t.Logf("%d: new roll file sink: %v", i, err)
			continue
		}
		fmt.Fprintf(sink, "[%s] %s\n", time.Now(), tt.url.String())
		sink.Close()
	}
}

func TestMain(m *testing.M) {
	os.RemoveAll("./testdata")
	m.Run()
	os.RemoveAll("./testdata")
}
