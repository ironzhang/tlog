package rollfile

import (
	"fmt"
	"io"
	"testing"
	"time"
)

func TestIsValidLayout(t *testing.T) {
	tests := []struct {
		layout string
		valid  bool
	}{
		{layout: "", valid: true},
		{layout: DayLayout, valid: true},
		{layout: HourLayout, valid: true},
		{layout: SecondLayout, valid: true},
		{layout: NanoLayout, valid: true},
		{layout: "20060103", valid: false},
	}
	for i, tt := range tests {
		valid := isValidLayout(tt.layout)
		if got, want := valid, tt.valid; got != want {
			t.Errorf("%d: valid: got %v, want %v", i, got, want)
			continue
		}
	}
}

func TestFileTick(t *testing.T) {
	f, err := Open("./testdata/test_file_tick/file.log")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()

	f.flushedAt = f.flushedAt.Add(-flushInterval)
	f.tick()

	f.maxSize = 1
	f.size = 1
	f.tick()

	f.maxSize = 0
	f.period = time.Second
	f.createdAt = f.createdAt.Add(-time.Second)
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

func TestFileOpen(t *testing.T) {
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

func TestFileSetLayout(t *testing.T) {
	tests := []struct {
		name   string
		layout string
	}{
		{name: "1.log", layout: ""},
		{name: "2.log", layout: DayLayout},
		{name: "3.log", layout: HourLayout},
		{name: "4.log", layout: SecondLayout},
		{name: "5.log", layout: NanoLayout},
	}
	for i, tt := range tests {
		f, err := Open("./testdata/test_file_set_layout/"+tt.name, SetLayout(tt.layout))
		if err != nil {
			t.Errorf("%d: open: %v", i, err)
			continue
		}
		PrintTestData(t, f, 1, "Hello, world\n")
		f.Close()
	}
}

func TestFileSetPeriod(t *testing.T) {
	f, err := Open("./testdata/test_file_period/file.log", SetPeriod(time.Second))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()

	PrintTestData(t, f, 1024, "Hello, world\n")
	//time.Sleep(2 * time.Second)
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
	f, err := Open("./testdata/test_file_print_create_log/file.log", PrintCreateLog())
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()

	PrintTestData(t, f, 1, "Hello, world\n")
}
