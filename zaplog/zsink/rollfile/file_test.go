package rollfile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestIsValidCutFormat(t *testing.T) {
	tests := []struct {
		format CutFormat
		valid  bool
	}{
		{format: "", valid: false},
		{format: SizeCut, valid: true},
		{format: HourCut, valid: true},
		{format: DayCut, valid: true},
		{format: "day", valid: false},
	}
	for i, tt := range tests {
		valid := isValidCutFormat(tt.format)
		if got, want := valid, tt.valid; got != want {
			t.Errorf("%d: valid: got %v, want %v", i, got, want)
			continue
		}
	}
}

func TestFileTick(t *testing.T) {
	f, err := Open("./testdata/test_file_tick/file.log", SetMaxSeq(2))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()

	f.flushedAt = f.flushedAt.Add(-flushInterval)
	f.tick()

	f.maxSize = 1
	f.size = 1
	f.tick()
}

func PrintTestData(t *testing.T, w io.Writer, n int, s string) {
	for i := 0; i < n; i++ {
		_, err := fmt.Fprint(w, s)
		if err != nil {
			t.Errorf("Fprint: %v", err)
		}
	}
}

func TestFileCheck(t *testing.T) {
	fileName := "./testdata/test_file_check/file.log"
	f, err := Open(fileName)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()

	data := []byte("hello, world")
	n, err := f.Write(data)
	if err != nil {
		t.Fatalf("write: %v", err)
	}
	if got, want := n, len(data); got != want {
		t.Errorf("n: got %d, want %d", got, want)
	}

	time.Sleep(2 * time.Second)

	// 删除文件及 link
	file := fileName
	link := filepath.Join(f.dir, f.name)

	os.Remove(file)
	os.Remove(link)

	time.Sleep(12 * time.Second)
	if !fileExist(fileName) {
		t.Fatalf("file is not exist %s", fileName)
	}
	if !fileExist(link) {
		t.Fatalf("link is not exist %s", link)
	}

	// 只删除 link
	os.Remove(link)
	time.Sleep(12 * time.Second)
	if !fileExist(link) {
		t.Fatalf("link is not exist %s", link)
	}

	n, err = f.Write(data)
	if err != nil {
		t.Fatalf("write: %v", err)
	}
	if got, want := n, len(data); got != want {
		t.Errorf("n: got %d, want %d", got, want)
	}
}

func TestFileWrite(t *testing.T) {
	f, err := Open("./testdata/test_file_write/file.log")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()

	data := []byte("hello, world")
	n, err := f.Write(data)
	if err != nil {
		t.Fatalf("write: %v", err)
	}
	if got, want := n, len(data); got != want {
		t.Errorf("n: got %d, want %d", got, want)
	}
}

func TestFileFlush(t *testing.T) {
	f, err := Open("./testdata/test_file_flush/file.log")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	PrintTestData(t, f, 1, "Hello, world\n")
	if err := f.Flush(); err != nil {
		t.Errorf("flush: %v", err)
	}
}

func TestFileSync(t *testing.T) {
	f, err := Open("./testdata/test_file_sync/file.log")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	PrintTestData(t, f, 1, "Hello, world\n")
	if err := f.Sync(); err != nil {
		t.Errorf("sync: %v", err)
	}
}

func TestFileClose(t *testing.T) {
	f, err := Open("./testdata/test_file_close/file.log")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	PrintTestData(t, f, 1, "Hello, world\n")
	if err := f.Close(); err != nil {
		t.Errorf("close: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Logf("close: %v", err)
	}
}

func TestFileSetCutFormat(t *testing.T) {
	tests := []struct {
		name   string
		format CutFormat
	}{
		{name: "1.log", format: SizeCut},
		{name: "2.log", format: HourCut},
		{name: "3.log", format: DayCut},
		//{name: "4.log", format: ""},
	}
	for i, tt := range tests {
		f, err := Open("./testdata/test_file_set_cut_format/"+tt.name, SetCutFormat(tt.format))
		if err != nil {
			t.Errorf("%d: open: %v", i, err)
			continue
		}
		PrintTestData(t, f, 1, "Hello, world\n")
		f.Close()
	}
}

func TestFileSetMaxSeq(t *testing.T) {
	f, err := Open("./testdata/test_file_set_max_seq/file.log", SetMaxSize(100), SetMaxSeq(3))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()

	PrintTestData(t, f, 1024, "Hello, world\n")
}

func TestFileSetMaxSize(t *testing.T) {
	f, err := Open("./testdata/test_file_size/file.log", SetMaxSize(1024))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()

	PrintTestData(t, f, 1024, "Hello, world\n")
}

func TestFilePrintCreateLog(t *testing.T) {
	PrintCreateLog = true
	f, err := Open("./testdata/test_file_print_create_log/file.log")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()

	PrintTestData(t, f, 1, "Hello, world\n")
}
