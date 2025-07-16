package main

import (
	"bufio"
	"fmt"
	"github.com/mikerowehl/gomal/pkg/reader"
	"os"
)

// Returned values are:
//
//	string - token/line
//	bool - eof, true means end of input
//	error - set to nil unless there's an error
func READ(scanner *bufio.Scanner) (reader.MalType, bool, error) {
	fmt.Print("user> ")

	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", false, fmt.Errorf("READ unable to get line: %w", err)
		} else {
			return "", true, nil
		}
	}
	val, err := reader.Read_str(scanner.Text())
	return val, false, err
}

func EVAL(v reader.MalType) reader.MalType {
	return v
}

func PRINT(v reader.MalType) {
	reader.Pr_str(v, true)
	fmt.Println()
}

func rep() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		line, eof, err := READ(scanner)
		if eof {
			fmt.Println()
			return
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
		} else {
			PRINT(EVAL(line))
		}
	}
}

func main() {
	rep()
}
