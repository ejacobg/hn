package sort

import (
	"bufio"
	"fmt"
	"os"
)

// writeString will attempt to write the given string to the given file.
func writeString(name string, flag int, s string) error {
	file, err := os.OpenFile(name, flag, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(s)
	if err != nil {
		return err
	}

	return nil
}

// processLines will call the given function on every line of a file.
func processLines(name string, processor func(string)) {
	file, err := os.Open(name)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		processor(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}
