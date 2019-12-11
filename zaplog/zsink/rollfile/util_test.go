package rollfile

import (
	"os"
)

func fileExist(name string) bool {
	_, err := os.Stat(name)
	if err != nil {
		return false
	}
	return true
}

//func TestFileName(t *testing.T) {
//	ts := time.Date(2019, 12, 6, 14, 1, 2, 3, time.Local)
//	tests := []struct {
//		prefix string
//		n      int
//		time   time.Time
//		layout string
//		name   string
//	}{
//		{
//			prefix: "debug.log",
//			n:      0,
//			time:   ts,
//			layout: "",
//			name:   "debug.log.0",
//		},
//		{
//			prefix: "debug.log",
//			n:      0,
//			time:   ts,
//			layout: "2006-01-02",
//			name:   "debug.log.0.2019-12-06",
//		},
//		{
//			prefix: "debug.log",
//			n:      0,
//			time:   ts,
//			layout: "2006-01-02-15",
//			name:   "debug.log.0.2019-12-06-14",
//		},
//		{
//			prefix: "debug.log",
//			n:      1,
//			time:   ts,
//			layout: "2006-01-02-15",
//			name:   "debug.log.1.2019-12-06-14",
//		},
//	}
//	for i, tt := range tests {
//		name := filename(tt.prefix, tt.n, tt.time, tt.layout)
//		if got, want := name, tt.name; got != want {
//			t.Errorf("%d: name: got %v, want %v", i, got, want)
//			continue
//		}
//		t.Logf("%d: name: got %v", i, name)
//	}
//}

//func TestIsSamePeriod(t *testing.T) {
//	tests := []struct {
//		t1     time.Time
//		t2     time.Time
//		period time.Duration
//		same   bool
//	}{
//		{
//			t1:     time.Date(2019, 12, 6, 14, 1, 2, 3, time.Local),
//			t2:     time.Date(2019, 12, 6, 15, 1, 2, 3, time.Local),
//			period: 24 * time.Hour,
//			same:   true,
//		},
//		{
//			t1:     time.Date(2019, 12, 6, 14, 1, 2, 3, time.Local),
//			t2:     time.Date(2019, 12, 7, 12, 1, 2, 3, time.Local),
//			period: 24 * time.Hour,
//			same:   false,
//		},
//		{
//			t1:     time.Date(2019, 12, 6, 14, 1, 2, 3, time.Local),
//			t2:     time.Date(2019, 12, 6, 15, 1, 2, 3, time.Local),
//			period: time.Hour,
//			same:   false,
//		},
//		{
//			t1:     time.Date(2019, 12, 6, 14, 1, 2, 3, time.Local),
//			t2:     time.Date(2019, 12, 6, 14, 59, 2, 3, time.Local),
//			period: time.Hour,
//			same:   true,
//		},
//		{
//			t1:     time.Date(2019, 12, 6, 14, 1, 2, 3, time.Local),
//			t2:     time.Date(2019, 12, 6, 14, 59, 2, 3, time.Local),
//			period: time.Hour / 2,
//			same:   false,
//		},
//		{
//			t1:     time.Date(2019, 12, 6, 14, 1, 2, 3, time.Local),
//			t2:     time.Date(2019, 12, 6, 14, 29, 2, 3, time.Local),
//			period: time.Hour / 2,
//			same:   true,
//		},
//	}
//	for i, tt := range tests {
//		same := isSamePeriod(tt.t1, tt.t2, tt.period)
//		if got, want := same, tt.same; got != want {
//			t.Errorf("%d: same: got %v, want %v", i, got, want)
//		}
//	}
//}

//func TestCreateDir(t *testing.T) {
//	dirs := []string{"a", "a/b", "b/c", "b", "a"}
//	for _, dir := range dirs {
//		if err := createDir(dir); err != nil {
//			t.Errorf("create dir %q: %v", dir, err)
//			continue
//		}
//		defer os.RemoveAll(dir)
//
//		if !fileExist(dir) {
//			t.Errorf("dir %q is not existed", dir)
//			continue
//		}
//	}
//}

//func TestCreateFile(t *testing.T) {
//	tests := []struct {
//		filename string
//		symlink  string
//	}{
//		{
//			filename: "debug.log.2019-12-06",
//			symlink:  "debug.log",
//		},
//	}
//	for _, tt := range tests {
//		f, err := createFile(tt.filename, tt.symlink)
//		if err != nil {
//			t.Errorf("create file %q: %v", tt.filename, err)
//			continue
//		}
//		if !fileExist(tt.filename) {
//			t.Errorf("file %q is not existed", tt.filename)
//			continue
//		}
//		if !fileExist(tt.symlink) {
//			t.Errorf("file %q is not existed", tt.symlink)
//			continue
//		}
//		fmt.Fprintf(f, "%s\n", tt.filename)
//
//		os.Remove(tt.symlink)
//		os.Remove(tt.filename)
//	}
//}
