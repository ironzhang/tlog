package rollfile

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTimeCutFileName(t *testing.T) {
	ts := time.Date(2019, 12, 6, 14, 1, 2, 3, time.Local)
	tests := []struct {
		base   string
		time   time.Time
		layout string
		name   string
	}{
		{
			base:   "debug.log",
			time:   ts,
			layout: dayLayout,
			name:   "debug.log.20191206",
		},
		{
			base:   "debug.log",
			time:   ts,
			layout: hourLayout,
			name:   "debug.log.2019120614",
		},
	}
	for i, tt := range tests {
		name := timeCutFileName(tt.base, tt.time, tt.layout)
		if got, want := name, tt.name; got != want {
			t.Errorf("%d: name: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: name: got %v", i, name)
	}
}

func TestSizeCutFileName(t *testing.T) {
	tests := []struct {
		pid  int
		base string
		seq  int
		name string
	}{
		{
			pid:  0,
			base: "debug.log",
			seq:  0,
			name: "debug.log.0",
		},
		{
			pid:  101,
			base: "debug.log",
			seq:  1,
			name: "debug.log.101.1",
		},
	}
	for i, tt := range tests {
		if tt.pid > 0 {
			LayoutPID = true
			pid = tt.pid
		} else {
			LayoutPID = false
		}
		name := sizeCutFileName(tt.base, tt.seq)
		if got, want := name, tt.name; got != want {
			t.Errorf("%d: name: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: name: got %v", i, name)
	}
	LayoutPID = false
	pid = os.Getpid()
}

func TestParseSeq(t *testing.T) {
	tests := []struct {
		name string
		seq  int
	}{
		{name: "debug.log.0", seq: 0},
		{name: "debug.log.1", seq: 1},
		{name: "debug.log.-1", seq: -1},
		{name: "debug.log", seq: 0},
		{name: "debug", seq: 0},
		{name: ".1", seq: 1},
		{name: "1", seq: 0},
		{name: "debug.log.20191206", seq: 20191206},
	}
	for i, tt := range tests {
		seq := parseSeq(tt.name)
		if got, want := seq, tt.seq; got != want {
			t.Errorf("%d: seq: got %v, want %v", i, got, want)
			continue
		} else {
			t.Logf("%d: seq: got %v", i, got)
		}
	}
}

func fileExist(name string) bool {
	_, err := os.Stat(name)
	if err != nil {
		return false
	}
	return true
}

func TestReadLink(t *testing.T) {
	dir := "testdata/test_read_link"
	os.MkdirAll(dir, os.ModePerm)

	tests := []struct {
		filename string
		symlink  string
	}{
		{filename: "f1", symlink: "s1"},
		{filename: "f2", symlink: "s2"},
	}
	for i, tt := range tests {
		fname := filepath.Join(dir, tt.filename)
		lname := filepath.Join(dir, tt.symlink)
		f, err := os.Create(fname)
		if err != nil {
			t.Fatalf("%d: os create: %v", i, err)
		}
		if err = os.Symlink(tt.filename, lname); err != nil {
			t.Fatalf("%d: os symlink: %v", i, err)
		}
		name, err := readLink(dir, tt.symlink)
		if err != nil {
			t.Fatalf("%d: read link: %v", i, err)
		}
		if got, want := name, tt.filename; got != want {
			t.Errorf("%d: filename: got %v, want %v", i, got, want)
		}
		f.Close()
	}
}

func TestReaLinkSeq(t *testing.T) {
	dir := "testdata/test_read_link_seq"
	os.MkdirAll(dir, os.ModePerm)

	tests := []struct {
		filename string
		symlink  string
		seq      int
	}{
		{filename: "", symlink: "debug.log", seq: 0},
		{filename: "debug.log.0", symlink: "debug.log", seq: 0},
		{filename: "info.log.1", symlink: "info.log", seq: 1},
	}
	for i, tt := range tests {
		if tt.filename != "" {
			fname := filepath.Join(dir, tt.filename)
			lname := filepath.Join(dir, tt.symlink)
			f, err := os.Create(fname)
			if err != nil {
				t.Fatalf("%d: os create: %v", i, err)
			}
			if err = os.Symlink(tt.filename, lname); err != nil {
				t.Fatalf("%d: os symlink: %v", i, err)
			}
			f.Close()
		}
		seq := readLinkSeq(dir, tt.symlink)
		if got, want := seq, tt.seq; got != want {
			t.Errorf("%d: seq: got %v, want %v", i, got, want)
		} else {
			t.Logf("%d: seq: got %v", i, got)
		}
	}
}

func TestCreateDir(t *testing.T) {
	dirs := []string{"testdata/test_create_dir/a", "testdata/test_create_dir/a", "testdata/test_create_dir/b/c"}
	for _, dir := range dirs {
		if err := createDir(dir); err != nil {
			t.Errorf("create dir %q: %v", dir, err)
			continue
		}

		if !fileExist(dir) {
			t.Errorf("dir %q is not existed", dir)
			continue
		}
	}
}

func TestCreateFile(t *testing.T) {
	dir := "testdata/test_create_file"
	os.MkdirAll(dir, os.ModePerm)
	os.MkdirAll(filepath.Join(dir, "log"), os.ModePerm)

	tests := []struct {
		filename string
		symlink  string
	}{
		{
			filename: "filename.0",
			symlink:  "symlink.0",
		},
		{
			filename: "log/filename.1",
			symlink:  "symlink.1",
		},
	}
	for i, tt := range tests {
		f, err := createFile(dir, tt.filename, tt.symlink)
		if err != nil {
			t.Errorf("%d: create file %q: %v", i, tt.filename, err)
			continue
		}

		if file := filepath.Join(dir, tt.filename); !fileExist(file) {
			t.Errorf("%d: file %q is not existed", i, file)
			continue
		}
		if file := filepath.Join(dir, tt.symlink); !fileExist(file) {
			t.Errorf("%d: file %q is not existed", i, file)
			continue
		}
		fmt.Fprintf(f, "%d: %s->%s\n", i, tt.symlink, tt.filename)
	}
}

func TestOpenFile(t *testing.T) {
	dir := "testdata/test_open_file"
	os.MkdirAll(dir, os.ModePerm)

	tests := []struct {
		filename string
		symlink  string
	}{
		{
			filename: "debug.log.0",
			symlink:  "debug.log",
		},
		{
			filename: "debug.log.1",
			symlink:  "debug.log",
		},
	}
	for i, tt := range tests {
		f, err := openFile(dir, tt.filename, tt.symlink)
		if err != nil {
			t.Errorf("%d: open file %q: %v", i, tt.filename, err)
			continue
		}

		if file := filepath.Join(dir, tt.filename); !fileExist(file) {
			t.Errorf("%d: file %q is not existed", i, file)
			continue
		}
		if file := filepath.Join(dir, tt.symlink); !fileExist(file) {
			t.Errorf("%d: file %q is not existed", i, file)
			continue
		}
		fmt.Fprintf(f, "%d: %s->%s\n", i, tt.symlink, tt.filename)
	}
}

func TestIsSamePeriod(t *testing.T) {
	tests := []struct {
		t1     time.Time
		t2     time.Time
		period time.Duration
		same   bool
	}{
		{
			t1:     time.Date(2019, 12, 6, 14, 1, 2, 3, time.Local),
			t2:     time.Date(2019, 12, 6, 15, 1, 2, 3, time.Local),
			period: 24 * time.Hour,
			same:   true,
		},
		{
			t1:     time.Date(2019, 12, 6, 14, 1, 2, 3, time.Local),
			t2:     time.Date(2019, 12, 7, 12, 1, 2, 3, time.Local),
			period: 24 * time.Hour,
			same:   false,
		},
		{
			t1:     time.Date(2019, 12, 6, 14, 1, 2, 3, time.Local),
			t2:     time.Date(2019, 12, 6, 15, 1, 2, 3, time.Local),
			period: time.Hour,
			same:   false,
		},
		{
			t1:     time.Date(2019, 12, 6, 14, 1, 2, 3, time.Local),
			t2:     time.Date(2019, 12, 6, 14, 59, 2, 3, time.Local),
			period: time.Hour,
			same:   true,
		},
		{
			t1:     time.Date(2019, 12, 6, 14, 1, 2, 3, time.Local),
			t2:     time.Date(2019, 12, 6, 14, 59, 2, 3, time.Local),
			period: time.Hour / 2,
			same:   false,
		},
		{
			t1:     time.Date(2019, 12, 6, 14, 1, 2, 3, time.Local),
			t2:     time.Date(2019, 12, 6, 14, 29, 2, 3, time.Local),
			period: time.Hour / 2,
			same:   true,
		},
	}
	for i, tt := range tests {
		same := isSamePeriod(tt.t1, tt.t2, tt.period)
		if got, want := same, tt.same; got != want {
			t.Errorf("%d: same: got %v, want %v", i, got, want)
		}
	}
}
