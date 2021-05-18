package rollfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func timeCutFileName(base string, t time.Time, layout string) string {
	if LayoutPID {
		return fmt.Sprintf("%s.%d.%s", base, pid, t.Format(layout))
	}
	return fmt.Sprintf("%s.%s", base, t.Format(layout))
}

func sizeCutFileName(base string, seq int) string {
	if LayoutPID {
		return fmt.Sprintf("%s.%d.%d", base, pid, seq)
	}
	return fmt.Sprintf("%s.%d", base, seq)
}

func parseSeq(name string) int {
	i := strings.LastIndex(name, ".")
	if i < 0 {
		return 0
	}
	n, err := strconv.Atoi(name[i+1:])
	if err != nil {
		return 0
	}
	return n
}

func readLink(dir, symlink string) (string, error) {
	link := filepath.Join(dir, symlink)
	return os.Readlink(link)
}

func readLinkSeq(dir, symlink string) int {
	file, err := readLink(dir, symlink)
	if err != nil {
		return 0
	}
	return parseSeq(file)
}

func createDir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

func createFile(dir, filename, symlink string) (f *os.File, err error) {
	file := filepath.Join(dir, filename)
	link := filepath.Join(dir, symlink)
	f, err = os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	os.Remove(link)
	os.Symlink(filename, link)
	return f, nil
}

func openFile(dir, filename, symlink string) (f *os.File, err error) {
	file := filepath.Join(dir, filename)
	link := filepath.Join(dir, symlink)
	f, err = os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	os.Remove(link)
	os.Symlink(filename, link)
	return f, nil
}

func isSamePeriod(t1, t2 time.Time, d time.Duration) bool {
	return t1.Truncate(d) == t2.Truncate(d)
}
