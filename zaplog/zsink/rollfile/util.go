package rollfile

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func fileName(prefix string, pid int, seq int, t time.Time, layout string) string {
	if layout == "" {
		return fmt.Sprintf("%s.%d.%d", prefix, pid, seq)
	}
	return fmt.Sprintf("%s.%d.%d.%s", prefix, pid, seq, t.Format(layout))
}

func createDir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

func createFile(filename, symlink string) (f *os.File, err error) {
	f, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	os.Remove(symlink)
	os.Symlink(filepath.Base(filename), symlink)
	return f, nil
}

func isSamePeriod(t1, t2 time.Time, d time.Duration) bool {
	return t1.Truncate(d) == t2.Truncate(d)
}
