package rollfile

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileName(t *testing.T) {
	ts := time.Date(2019, 12, 6, 14, 1, 2, 3, time.Local)
	tests := []struct {
		prefix string
		pid    int
		seq    int
		time   time.Time
		layout string
		name   string
	}{
		{
			prefix: "debug.log",
			pid:    1,
			seq:    2,
			time:   ts,
			layout: "",
			name:   "debug.log.1.2",
		},
		{
			prefix: "debug.log",
			pid:    0,
			seq:    0,
			time:   ts,
			layout: "",
			name:   "debug.log.0.0",
		},
		{
			prefix: "debug.log",
			pid:    0,
			seq:    0,
			time:   ts,
			layout: DayLayout,
			name:   "debug.log.0.0.20191206",
		},
		{
			prefix: "debug.log",
			pid:    0,
			seq:    0,
			time:   ts,
			layout: HourLayout,
			name:   "debug.log.0.0.2019120614",
		},
		{
			prefix: "debug.log",
			pid:    0,
			seq:    0,
			time:   ts,
			layout: SecondLayout,
			name:   "debug.log.0.0.20191206-140102",
		},
		{
			prefix: "debug.log",
			pid:    0,
			seq:    0,
			time:   ts,
			layout: NanoLayout,
			name:   "debug.log.0.0.20191206-140102.000000003",
		},
	}
	for i, tt := range tests {
		name := fileName(tt.prefix, tt.pid, tt.seq, tt.time, tt.layout)
		if got, want := name, tt.name; got != want {
			t.Errorf("%d: name: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: name: got %v", i, name)
	}
}

func fileExist(name string) bool {
	_, err := os.Stat(name)
	if err != nil {
		return false
	}
	return true
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
