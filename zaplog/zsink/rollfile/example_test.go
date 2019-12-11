package rollfile

import (
	"fmt"
	"os"
)

func ExampleFile() {
	os.RemoveAll("./testdata")

	file, err := Open("./testdata/log/debug.log", SetLayout(HourLayout), SetMaxSize(80))
	if err != nil {
		fmt.Printf("open: %v", err)
		return
	}
	defer file.Close()

	n := 10
	for i := 0; i < n; i++ {
		fmt.Fprintf(file, "hello\n")
	}

	// output:
}
