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

func TestSuffixToLayout(t *testing.T) {
	tests := []struct {
		suffix string
		layout string
		err    string
	}{
		{suffix: "d", layout: rollfile.DayLayout},
		{suffix: "day", layout: rollfile.DayLayout},
		{suffix: "h", layout: rollfile.HourLayout},
		{suffix: "hour", layout: rollfile.HourLayout},
		{suffix: "s", layout: rollfile.SecondLayout},
		{suffix: "second", layout: rollfile.SecondLayout},
		{suffix: "n", layout: rollfile.NanoLayout},
		{suffix: "nano", layout: rollfile.NanoLayout},
		{suffix: "Day", layout: rollfile.DayLayout},
		{suffix: "DAY", layout: rollfile.DayLayout},
		{suffix: "days", err: "unknown suffix pattern"},
	}
	for i, tt := range tests {
		layout, err := suffixToLayout(tt.suffix)
		if !matchError(t, err, tt.err) {
			t.Errorf("%d: match error: got %v, want %v", i, err, tt.err)
			continue
		}
		if err != nil {
			t.Logf("%d: suffix to layout: %v", i, err)
			continue
		}
		if got, want := layout, tt.layout; got != want {
			t.Errorf("%d: layout: got %v, want %v", i, got, want)
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
		{url: ParseTestURL(t, "rfile://workdir/testdata/log/b.log?suffix=hour")},
		{url: ParseTestURL(t, "rfile://workdir/testdata/log/f.log?suffix=year"), err: "unknown suffix pattern"},
		{url: ParseTestURL(t, "rfile://workdir/testdata/log/c.log?period=1h")},
		{url: ParseTestURL(t, "rfile://workdir/testdata/log/c.log?period=1hour"), err: "unknown unit"},
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
