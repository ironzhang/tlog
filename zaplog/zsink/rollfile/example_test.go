package rollfile

import (
	"fmt"
	"os"
	"testing"
)

func ExampleFile() {
	file, err := Open("./testdata/example_file/debug.log", SetLayout(HourLayout), SetMaxSize(80))
	if err != nil {
		fmt.Printf("open: %v", err)
		return
	}
	defer file.Close()

	n := 1
	for i := 0; i < n; i++ {
		fmt.Fprintf(file, "hello\n")
	}

	//time.Sleep(10 * time.Second)

	// output:
}

func TestMain(m *testing.M) {
	os.RemoveAll("./testdata")
	m.Run()
	os.RemoveAll("./testdata")
}
