package rollfile

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func ExampleFile() {
	file, err := Open("./testdata/example_file/debug.log", SetMaxSize(10), SetMaxSeq(10))
	//file, err := Open("./testdata/example_file/debug.log", SetCutFormat(HourCut), SetMaxSize(10), SetMaxSeq(10))
	if err != nil {
		fmt.Printf("open: %v", err)
		return
	}
	defer file.Close()

	n := 5
	for i := 0; i < n; i++ {
		fmt.Fprintf(file, "hello\n")
	}

	//time.Sleep(10 * time.Second)

	// output:
}

func TestMain(m *testing.M) {
	// 让单测时间变短
	flushInterval, checkInterval = time.Second, time.Second

	os.RemoveAll("./testdata")
	m.Run()
	os.RemoveAll("./testdata")
}
