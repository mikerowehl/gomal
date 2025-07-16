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

func APPLY(v reader.MalType) (reader.MalType, error) {
	if l, ok := v.(reader.MalList); !ok {
		fmt.Println("Applying something that isn't a list")
	} else {
		a, ok := l[0].(reader.MalFunc)
		if !ok {
			return nil, fmt.Errorf("Error converting apply function")
		}
		return a(l[1:])
	}
	return nil, fmt.Errorf("Error applying")
}

func listCast(raw reader.MalType, min int) (reader.MalList, error) {
	l, ok := raw.(reader.MalList)
	if !ok {
		return nil, fmt.Errorf("listCast: expected list: %v", raw)
	}
	if len(l) < min {
		return nil, fmt.Errorf("listCast: expected at least %d elements, got %d", min, len(l))
	}
	return l, nil
}

var add = func(a reader.MalType) (reader.MalType, error) {
	l, err := listCast(a, 2)
	if err != nil {
		return nil, err
	}
	i1, ok := l[0].(int)
	if !ok {
		return nil, fmt.Errorf("add: non int first arg to add: %v", l[0])
	}
	i2, ok := l[1].(int)
	if !ok {
		return nil, fmt.Errorf("add: non int second arg to add: %v", l[1])
	}
	return reader.MalType(i1 + i2), nil
}

var sub = func(a reader.MalType) (reader.MalType, error) {
	l, err := listCast(a, 2)
	if err != nil {
		return nil, err
	}
	i1, ok := l[0].(int)
	if !ok {
		return nil, fmt.Errorf("sub: non int first arg: %v", l[0])
	}
	i2, ok := l[1].(int)
	if !ok {
		return nil, fmt.Errorf("sub: non int second arg: %v", l[1])
	}
	return reader.MalType(i1 - i2), nil
}

var mul = func(a reader.MalType) (reader.MalType, error) {
	l, err := listCast(a, 2)
	if err != nil {
		return nil, err
	}
	i1, ok := l[0].(int)
	if !ok {
		return nil, fmt.Errorf("mul: non int first arg: %v", l[0])
	}
	i2, ok := l[1].(int)
	if !ok {
		return nil, fmt.Errorf("mul: non int second arg: %v", l[1])
	}
	return reader.MalType(i1 * i2), nil
}

var div = func(a reader.MalType) (reader.MalType, error) {
	l, err := listCast(a, 2)
	if err != nil {
		return nil, err
	}
	i1, ok := l[0].(int)
	if !ok {
		return nil, fmt.Errorf("div: non int first arg: %v", l[0])
	}
	i2, ok := l[1].(int)
	if !ok {
		return nil, fmt.Errorf("div: non int second arg: %v", l[1])
	}
	return reader.MalType(i1 / i2), nil
}
var env = map[string]reader.MalType{
	"+": reader.MalFunc(add),
	"-": reader.MalFunc(sub),
	"*": reader.MalFunc(mul),
	"/": reader.MalFunc(div),
}

func EVAL(v reader.MalType, env map[string]reader.MalType) reader.MalType {
	switch t := v.(type) {
	case reader.MalSymbol:
		for name, entry := range env {
			if name == string(t) {
				return entry
			}
		}
	case reader.MalList:
		if len(t) == 0 {
			return v
		}
		evaled := reader.MalList{}
		for _, entry := range t {
			n := EVAL(entry, env)
			evaled = append(evaled, n)
		}
		app, err := APPLY(evaled)
		if err != nil {
			fmt.Println("Error returned from apply")
			return nil
		}
		return app
	case reader.MalVector:
		evaled := reader.MalVector{}
		for _, entry := range t {
			n := EVAL(entry, env)
			evaled = append(evaled, n)
		}
		return evaled
	case reader.MalHashmap:
		evaled := reader.MalHashmap{}
		for i, entry := range t {
			if (i % 2) == 1 {
				n := EVAL(entry, env)
				evaled = append(evaled, n)
			} else {
				evaled = append(evaled, entry)
			}
		}
		return evaled
	}
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
			PRINT(EVAL(line, env))
		}
	}
}

func main() {
	rep()
}
