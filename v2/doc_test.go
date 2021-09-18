package errors_test

import (
	"fmt"
	"os"

	"github.com/payfazz/go-errors/v2"
)

func readFile() (string, error) {
	data, err := os.ReadFile("InvalidFile.txt")
	if err != nil {
		return "", errors.Trace(err)
	}

	return string(data), nil
}

func doSomething() error {
	data, err := readFile()
	if err != nil {
		return errors.Trace(err)
	}

	fmt.Println(data)

	return nil
}

func Example() {
	if err := errors.Catch(doSomething); err != nil {
		for _, loc := range errors.StackTrace(err) {
			fmt.Fprintln(os.Stderr, loc.String())
		}
	}
}
